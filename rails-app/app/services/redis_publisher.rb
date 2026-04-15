require "redis"
require "json"

class RedisPublisher
  def self.publish_ticket(ticket)
    redis = Redis.new

    payload = {
      id: ticket.id,
      title: ticket.title,
      description: ticket.description
    }

    redis.publish("tickets", payload.to_json)
  end
end