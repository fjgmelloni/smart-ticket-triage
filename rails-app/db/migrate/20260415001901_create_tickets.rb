class CreateTickets < ActiveRecord::Migration[7.1]
  def change
    create_table :tickets do |t|
      t.string :title
      t.text :description
      t.string :status
      t.string :category
      t.string :priority
      t.text :ai_summary

      t.timestamps
    end
  end
end
