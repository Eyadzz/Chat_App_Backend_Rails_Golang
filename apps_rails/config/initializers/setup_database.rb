# config/initializers/setup_database.rb

# Check if the database exists, and if not, create it
ActiveRecord::Tasks::DatabaseTasks.create_current
# Check if migrations need to be applied
ActiveRecord::MigrationContext.new('db/migrate').migrate unless ActiveRecord::Migrator.current_version.zero?
