class AppsService
  def self.create_app(name)
    app = App.new(name: name, token: SecureRandom.uuid)
    app.save
    app
  end

  def self.update_app(token, name)
    app = App.find_by(token: token)

    if app.nil?
      return nil
    end

    app.update(name: name)
    app
  end

  def self.get_all
    App.all
  end

  def self.get_app(token)
    App.find_by(token: token)
  end
end