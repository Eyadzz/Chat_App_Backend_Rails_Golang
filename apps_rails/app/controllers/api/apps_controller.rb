module Api
  class AppsController < ApplicationController

    def index
      render json: AppSerializer.new(AppsService.get_all).serialized_json
    end

    def show
      render json: AppSerializer.new(AppsService.get_app(params[:token])).serialized_json
    end

    def create
      app = AppsService.create_app(app_params[:name])

      if app.persisted?
        render json: { token: app.token }, status: 201
      else
        render json: { error: app.errors.messages }, status: 422
      end
    end

    def update
      app = AppsService.update_app(params[:token], app_params[:name])

      if app.nil?
        render json: { error: "App not found" }, status: 404
      else
        render json: { status: "App updated" }, status: 200
      end
    end

    private def app_params
      params.require(:app).permit(:name)
    end
  end
end