#!/usr/bin/ruby
# frozen_string_literal: true

require 'net/http'
require 'json'
require 'pathname'
require 'time'
require 'logger'

class LeaderElector
  KUBE_API_SERVER = 'https://kubernetes.default.svc'

  def initialize
    @logger = Logger.new($stdout)
    @logger.level = Logger::DEBUG

    @ttl = 10
    @pause = 5

    # setup k8s info
    @directory = '/run/secrets/kubernetes.io/serviceaccount'
    @namespace = 'proxysql-test'
    @lease_lock = 'proxysql-leader-lease'
  end

  # get the lease for the leader lock, if it exists
  def get_lease
    uri = URI("#{KUBE_API_SERVER}/apis/coordination.k8s.io/v1/namespaces/#{@namespace}/leases/#{@lease_lock}")
    token = File.read(Pathname.new("#{@directory}/token").realpath)

    http = Net::HTTP.new(uri.host, uri.port)
    http.use_ssl = true
    http.ca_file = Pathname.new("#{@directory}/ca.crt").realpath.to_s

    request = Net::HTTP::Get.new(uri)
    request['Authorization'] = "Bearer #{token}"

    response = http.request(request)

    unless response.is_a?(Net::HTTPSuccess)
      @logger.error "get_lease - request failed with code: #{response.code}"
      @logger.error "Error message: #{response.body}"

      return false
    end

    if response.code.to_i != 200
      @logger.debug "Lease doesn't exist"
      return nil
    end

    @lease = JSON.parse(response.body)
    @lease
  end

  # if the lease lock does not exist, create a new one. confers leadership.
  def create_lease
    uri = URI("#{KUBE_API_SERVER}/apis/coordination.k8s.io/v1/namespaces/#{@namespace}/leases")
    token = File.read(Pathname.new("#{@directory}/token").realpath)

    lease_data = {
      metadata: {
        name: @lease_lock
      },
      spec: {
        holderIdentity: ENV.fetch('HOSTNAME', 'idk'),
        leaseDurationSeconds: 10
      }
    }

    http = Net::HTTP.new(uri.host, uri.port)
    http.use_ssl = true
    http.ca_file = Pathname.new("#{@directory}/ca.crt").realpath.to_s

    request = Net::HTTP::Post.new(uri)
    request['Authorization'] = "Bearer #{token}"
    request.content_type = 'application/json'
    request.body = lease_data.to_json

    response = http.request(request)

    unless response.is_a?(Net::HTTPSuccess)
      # leader is already taken
      return false if response.code.to_i == 409

      @logger.error "create_lease - request failed with code: #{response.code}"
      @logger.error "Error message: #{response.body}"

      return false
    end

    if response.code.to_i != 200
      @logger.warn "create_lease - something weird happened: #{response.code}"
      @logger.warn "response: #{response.body}"
    else
      @logger.debug 'Lease doesnt exist, creating a new one'
    end
  end

  # renew the lease lock, keeping leadership
  def renew_lease
    uri = URI("#{KUBE_API_SERVER}/apis/coordination.k8s.io/v1/namespaces/#{@namespace}/leases/#{@lease_lock}")
    token = File.read(Pathname.new("#{@directory}/token").realpath)

    http = Net::HTTP.new(uri.host, uri.port)
    http.use_ssl = true
    http.ca_file = Pathname.new("#{@directory}/ca.crt").realpath.to_s

    new_lease = @lease.dup

    new_lease['spec']['renewTime'] = Time.now.strftime('%Y-%m-%dT%H:%M:%S.%6N%:z') # Time.now.utc.iso8601
    new_lease['spec']['holderIdentity'] = ENV.fetch('HOSTNAME')

    # FIXME: log "took master" or whatever, if holderIdentity is different from hostname

    request = Net::HTTP::Put.new(uri)
    request['Content-Type'] = 'application/json'
    request['Authorization'] = "Bearer #{token}"
    request.body = new_lease.to_json

    response = http.request(request)

    puts @lease['spec']['holderIdentity']
    puts new_lease['spec']['holderIdentity']

    if response.is_a?(Net::HTTPSuccess)
      if @lease != new_lease
        @logger.debug "#{ENV.fetch('HOSTNAME')} took lease pver from #{@lease['sped']['holderIdentity']}"
      else
        @logger.debug "Renewed lease for #{ENV.fetch('HOSTNAME')}"
      end

      @lease = new_lease
    else
      @logger.error "renew_lease - request failed with code: #{response.code}"
      @logger.error "Error message: #{response.body}"

      false
    end
  end


  def lease_owner?
    holder_identity = @lease['spec']['holderIdentity']
    holder_identity == ENV.fetch('HOSTNAME')
  end


  def lease_expired?
    ttl = @lease['spec'].fetch('leaseDurationSeconds', @ttl) # get ttl or default to 30

    # get the last renewTime from the lease, OR use creationTimestamp of it's null (as a new lease would be)
    renewed_at = @lease['spec'].fetch('renewTime', nil) || @lease['metadata']['creationTimestamp']

    renewed_at = Time.parse(renewed_at)

    expiration_time = renewed_at + ttl

    Time.now > expiration_time
  end

  def run
    loop do
      get_lease

      if @lease.nil?
        create_lease
      elsif lease_owner? || lease_expired?
        renew_lease
      else
        @logger.debug 'Nothing to do, sleeping.'
      end

      sleep(@pause)
    end
  end
end

le = LeaderElector.new
le.run
