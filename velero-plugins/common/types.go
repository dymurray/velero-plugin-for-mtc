package common

const (
	MigrationRegistry string = "openshift.io/migration-registry"
	// distinction for B/R and migration
	MigrationApplicationLabelKey   string = "app.kubernetes.io/part-of"
	MigrationApplicationLabelValue string = "openshift-migration"
	MigMigrationLabelKey           string = "migration.openshift.io/migrated-by-migmigration"
	MigPlanLabelKey                string = "migration.openshift.io/migrated-by-migplan"
)
