class AppSerializer
  include FastJsonapi::ObjectSerializer
  attributes :name, :token
end
