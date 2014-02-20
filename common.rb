require 'uri'
require 'cgi'
require 'open-uri'
require 'json'
require 'pathname'

# The pattern used to check a commit message contains a Pivotal Tracker story ID
TAG_PATTERN = /
  \[
    (?:
      (?:
       (?:complete[sd]?|(?:finish|fix)(?:e[sd])?)\s+
      )?\#\d{4,}
    |
      \#?no[ ]?story
    )
  \]
/ix

TRACKER_STORIES_CACHE_FILEPATH = Pathname.new(
  File.expand_path('../.cache/tracker_stories', __FILE__))

module Colors
  def red(string)
    "\e[31m#{string}\e[m"
  end
end

class Selecta
  # Get selecta here: https://github.com/garybernhardt/selecta
  def self.available?
    location = `which selecta`.chomp
    File.executable?(location)
  end

  def self.prompt(list)
    IO.popen('selecta 2> /dev/null', 'r+') do |io|
      io.puts list
      io.close_write
      io.gets.strip
    end
  rescue Interrupt
    nil
  end
end

class PivotalAPI
  PIVOTAL_API_URL = 'https://www.pivotaltracker.com/services/v5/'

  def initialize(api_token)
    raise 'No Pivotal Tracker API token available.' if api_token.empty?
    @token = api_token
  end

  def my_active_stories
    filter = "state:started mywork:#{me['username']}"
    my_project_ids.flat_map do |project_id|
      stories_for_project(project_id, filter: filter)
    end
  end

  private

  def stories_for_project(project_id, params = {})
    api_request("projects/#{project_id}/stories", params)
  end

  def me
    @me ||= api_request('me')
  end

  def my_project_ids
    me['projects'].map { |project| project['project_id'] }
  end

  def api_request(path, params = {})
    uri = URI.parse(PIVOTAL_API_URL + path)
    uri.query = uri_query_from_hash(params)
    raw = open(uri.to_s, 'X-TrackerToken' => @token).read
    raise 'empty response' if raw.empty?
    JSON.load(raw)
  end

  def uri_query_from_hash(hash)
    pairs = hash.map do |key, value|
      '%s=%s' % [CGI.escape(key.to_s), CGI.escape(value.to_s)]
    end
    pairs.join('&')
  end
end

