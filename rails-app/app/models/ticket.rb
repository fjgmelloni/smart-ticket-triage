class Ticket < ApplicationRecord
  after_create :send_to_queue

  def send_to_queue
    RedisPublisher.publish_ticket(self)
  end
end