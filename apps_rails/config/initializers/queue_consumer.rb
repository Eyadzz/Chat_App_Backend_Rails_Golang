# config/initializers/message_consumer.rb

def start_queue_consumer
  Thread.new do
    connection = Bunny.new(ENV['RABBITMQ_URL'])
    connection.start

    channel = connection.create_channel
    queue = channel.queue('chats_count')

    begin
      puts ' [*] Waiting for messages. To exit press CTRL+C'
      queue.subscribe(block: true) do |_delivery_info, _properties, body|
        data = JSON.parse(body)

        application_token = data['ApplicationToken']
        number = data['Number']

        app = App.find_by(token: application_token)
        if app
          app.chats_count = number
          app.save
          puts " [x] Updated chats_count for App with token #{application_token} to #{app.chats_count}"
        else
          puts " [x] App with token #{application_token} not found"
        end
      end
    rescue Interrupt => _
      connection.close
      exit(0)
    end
  end
end

# Start the message consumer when Rails boots up
start_queue_consumer if defined?(Rails::Server)
