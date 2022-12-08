require "sinatra"

# Authentication token
$token = "test"

# List of all availble resources
$resources = {
  "id1" => "postgres://localhost:5432/db1?sslmode=disable",
  "id2" => "postgres://localhost:5432/db2?sslmode=disable",
  "id3" => "postgres://localhost:5432/db3?sslmode=disable"
}

helpers do
  def error(code, message)
    halt(code, JSON.dump(error: message))
  end
end

before do
  content_type :json
end

post "/" do
  req = JSON.load(request.body) || {}

  unless req["resource"]
    halt 404, "Resource ID required"
  end

  # Check the resource
  resource = $resources[req["resource"]]
  if !resource
    halt 404, "Invalid resource ID"
  end

  # Return connection credentials
  JSON.dump(
    database_url: resource
  )
end
