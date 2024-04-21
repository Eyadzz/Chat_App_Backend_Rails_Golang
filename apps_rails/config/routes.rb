Rails.application.routes.draw do

  namespace :api do
      resources :apps, param: :token, only: [:index, :create, :show, :update]
    end
end
