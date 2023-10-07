#!/usr/bin/env ruby
# frozen_string_literal: true

require 'json'
require 'net/http'
require 'pathname'
require 'uri'

directory = '/run/secrets/kubernetes.io/serviceaccount'

namespace = File.read(Pathname.new("#{directory}/namespace").realpath)
token = File.read(Pathname.new("#{directory}/token").realpath)

uri = URI("https://kubernetes.default.svc/api/v1/namespaces/#{namespace}/pods")

http = Net::HTTP.new(uri.host, uri.port)
http.use_ssl = true
http.ca_file = Pathname.new("#{directory}/ca.crt").realpath.to_s

request = Net::HTTP::Get.new(uri)
request['Authorization'] = "Bearer #{token}" # Set the Authorization header

response = http.request(request)

unless response.is_a?(Net::HTTPSuccess)
  puts "Request failed with code: #{response.code}"
  puts "Error message: #{response.body}"
  exit 1
end

data = JSON.parse(response.body)
pods = data['items']

puts 'DELETE FROM proxysql_servers;'

filtered_pods = pods.reject { |p| p['metadata']['labels']['app'] == 'proxysql-admin' }
filtered_pods.each do |pod|
  pod_ip = pod['status']['podIP']
  pod_name = pod['metadata']['name']

  sql_statement = "INSERT INTO proxysql_servers VALUES ('#{pod_ip}', 6032, 0, '#{pod_name}');"
  puts sql_statement
end

puts 'LOAD PROXYSQL SERVERS TO RUNTIME;'
puts 'LOAD MYSQL SERVERS TO RUNTIME;'
puts 'LOAD MYSQL USERS TO RUNTIME;'
puts 'LOAD MYSQL QUERY RULES TO RUNTIME;'
