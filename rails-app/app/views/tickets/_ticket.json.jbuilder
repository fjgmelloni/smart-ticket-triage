json.extract! ticket, :id, :title, :description, :status, :category, :priority, :ai_summary, :created_at, :updated_at
json.url ticket_url(ticket, format: :json)
