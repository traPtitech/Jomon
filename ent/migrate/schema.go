// Code generated by entc, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// CommentsColumns holds the columns for the "comments" table.
	CommentsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID},
		{Name: "comment", Type: field.TypeString},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "deleted_at", Type: field.TypeTime, Nullable: true},
		{Name: "request_comment", Type: field.TypeUUID, Nullable: true},
	}
	// CommentsTable holds the schema information for the "comments" table.
	CommentsTable = &schema.Table{
		Name:       "comments",
		Columns:    CommentsColumns,
		PrimaryKey: []*schema.Column{CommentsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "comments_requests_comment",
				Columns:    []*schema.Column{CommentsColumns[5]},
				RefColumns: []*schema.Column{RequestsColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// FilesColumns holds the columns for the "files" table.
	FilesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID},
		{Name: "name", Type: field.TypeString},
		{Name: "mime_type", Type: field.TypeString},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "deleted_at", Type: field.TypeTime, Nullable: true},
		{Name: "request_file", Type: field.TypeUUID, Nullable: true},
	}
	// FilesTable holds the schema information for the "files" table.
	FilesTable = &schema.Table{
		Name:       "files",
		Columns:    FilesColumns,
		PrimaryKey: []*schema.Column{FilesColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "files_requests_file",
				Columns:    []*schema.Column{FilesColumns[5]},
				RefColumns: []*schema.Column{RequestsColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// GroupsColumns holds the columns for the "groups" table.
	GroupsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID},
		{Name: "name", Type: field.TypeString},
		{Name: "description", Type: field.TypeString},
		{Name: "budget", Type: field.TypeInt, Nullable: true},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "deleted_at", Type: field.TypeTime, Nullable: true},
	}
	// GroupsTable holds the schema information for the "groups" table.
	GroupsTable = &schema.Table{
		Name:        "groups",
		Columns:     GroupsColumns,
		PrimaryKey:  []*schema.Column{GroupsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{},
	}
	// GroupBudgetsColumns holds the columns for the "group_budgets" table.
	GroupBudgetsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID},
		{Name: "amount", Type: field.TypeInt},
		{Name: "comment", Type: field.TypeString, Nullable: true, Size: 2147483647},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "group_group_budget", Type: field.TypeUUID, Nullable: true},
	}
	// GroupBudgetsTable holds the schema information for the "group_budgets" table.
	GroupBudgetsTable = &schema.Table{
		Name:       "group_budgets",
		Columns:    GroupBudgetsColumns,
		PrimaryKey: []*schema.Column{GroupBudgetsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "group_budgets_groups_group_budget",
				Columns:    []*schema.Column{GroupBudgetsColumns[4]},
				RefColumns: []*schema.Column{GroupsColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// RequestsColumns holds the columns for the "requests" table.
	RequestsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID},
		{Name: "amount", Type: field.TypeInt},
		{Name: "title", Type: field.TypeString},
		{Name: "content", Type: field.TypeString},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "group_request", Type: field.TypeUUID, Nullable: true},
	}
	// RequestsTable holds the schema information for the "requests" table.
	RequestsTable = &schema.Table{
		Name:       "requests",
		Columns:    RequestsColumns,
		PrimaryKey: []*schema.Column{RequestsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "requests_groups_request",
				Columns:    []*schema.Column{RequestsColumns[6]},
				RefColumns: []*schema.Column{GroupsColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// RequestStatusColumns holds the columns for the "request_status" table.
	RequestStatusColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID},
		{Name: "status", Type: field.TypeEnum, Enums: []string{"submitted", "fix_required", "accepted", "completed", "rejected"}, Default: "submitted"},
		{Name: "reason", Type: field.TypeString},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "request_status", Type: field.TypeUUID, Nullable: true},
	}
	// RequestStatusTable holds the schema information for the "request_status" table.
	RequestStatusTable = &schema.Table{
		Name:       "request_status",
		Columns:    RequestStatusColumns,
		PrimaryKey: []*schema.Column{RequestStatusColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "request_status_requests_status",
				Columns:    []*schema.Column{RequestStatusColumns[4]},
				RefColumns: []*schema.Column{RequestsColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// RequestTargetsColumns holds the columns for the "request_targets" table.
	RequestTargetsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID},
		{Name: "target", Type: field.TypeString},
		{Name: "paid_at", Type: field.TypeTime, Nullable: true},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "request_target", Type: field.TypeUUID, Nullable: true},
	}
	// RequestTargetsTable holds the schema information for the "request_targets" table.
	RequestTargetsTable = &schema.Table{
		Name:       "request_targets",
		Columns:    RequestTargetsColumns,
		PrimaryKey: []*schema.Column{RequestTargetsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "request_targets_requests_target",
				Columns:    []*schema.Column{RequestTargetsColumns[4]},
				RefColumns: []*schema.Column{RequestsColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// TagsColumns holds the columns for the "tags" table.
	TagsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID},
		{Name: "name", Type: field.TypeString},
		{Name: "description", Type: field.TypeString},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "deleted_at", Type: field.TypeTime, Nullable: true},
		{Name: "request_tag", Type: field.TypeUUID, Nullable: true},
		{Name: "transaction_tag", Type: field.TypeUUID, Nullable: true},
	}
	// TagsTable holds the schema information for the "tags" table.
	TagsTable = &schema.Table{
		Name:       "tags",
		Columns:    TagsColumns,
		PrimaryKey: []*schema.Column{TagsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "tags_requests_tag",
				Columns:    []*schema.Column{TagsColumns[6]},
				RefColumns: []*schema.Column{RequestsColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:     "tags_transactions_tag",
				Columns:    []*schema.Column{TagsColumns[7]},
				RefColumns: []*schema.Column{TransactionsColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// TransactionsColumns holds the columns for the "transactions" table.
	TransactionsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "group_budget_transaction", Type: field.TypeUUID, Unique: true, Nullable: true},
		{Name: "request_transaction", Type: field.TypeUUID, Nullable: true},
	}
	// TransactionsTable holds the schema information for the "transactions" table.
	TransactionsTable = &schema.Table{
		Name:       "transactions",
		Columns:    TransactionsColumns,
		PrimaryKey: []*schema.Column{TransactionsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "transactions_group_budgets_transaction",
				Columns:    []*schema.Column{TransactionsColumns[2]},
				RefColumns: []*schema.Column{GroupBudgetsColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:     "transactions_requests_transaction",
				Columns:    []*schema.Column{TransactionsColumns[3]},
				RefColumns: []*schema.Column{RequestsColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// TransactionDetailsColumns holds the columns for the "transaction_details" table.
	TransactionDetailsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID},
		{Name: "amount", Type: field.TypeInt, Default: 0},
		{Name: "target", Type: field.TypeString, Default: ""},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "transaction_detail", Type: field.TypeUUID, Unique: true, Nullable: true},
	}
	// TransactionDetailsTable holds the schema information for the "transaction_details" table.
	TransactionDetailsTable = &schema.Table{
		Name:       "transaction_details",
		Columns:    TransactionDetailsColumns,
		PrimaryKey: []*schema.Column{TransactionDetailsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "transaction_details_transactions_detail",
				Columns:    []*schema.Column{TransactionDetailsColumns[5]},
				RefColumns: []*schema.Column{TransactionsColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// UsersColumns holds the columns for the "users" table.
	UsersColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID},
		{Name: "name", Type: field.TypeString, Unique: true},
		{Name: "display_name", Type: field.TypeString},
		{Name: "admin", Type: field.TypeBool, Default: false},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "deleted_at", Type: field.TypeTime, Nullable: true},
		{Name: "comment_user", Type: field.TypeUUID, Unique: true, Nullable: true},
		{Name: "request_user", Type: field.TypeUUID, Unique: true, Nullable: true},
		{Name: "request_status_user", Type: field.TypeUUID, Unique: true, Nullable: true},
	}
	// UsersTable holds the schema information for the "users" table.
	UsersTable = &schema.Table{
		Name:       "users",
		Columns:    UsersColumns,
		PrimaryKey: []*schema.Column{UsersColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "users_comments_user",
				Columns:    []*schema.Column{UsersColumns[7]},
				RefColumns: []*schema.Column{CommentsColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:     "users_requests_user",
				Columns:    []*schema.Column{UsersColumns[8]},
				RefColumns: []*schema.Column{RequestsColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:     "users_request_status_user",
				Columns:    []*schema.Column{UsersColumns[9]},
				RefColumns: []*schema.Column{RequestStatusColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// GroupUserColumns holds the columns for the "group_user" table.
	GroupUserColumns = []*schema.Column{
		{Name: "group_id", Type: field.TypeUUID},
		{Name: "user_id", Type: field.TypeUUID},
	}
	// GroupUserTable holds the schema information for the "group_user" table.
	GroupUserTable = &schema.Table{
		Name:       "group_user",
		Columns:    GroupUserColumns,
		PrimaryKey: []*schema.Column{GroupUserColumns[0], GroupUserColumns[1]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "group_user_group_id",
				Columns:    []*schema.Column{GroupUserColumns[0]},
				RefColumns: []*schema.Column{GroupsColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "group_user_user_id",
				Columns:    []*schema.Column{GroupUserColumns[1]},
				RefColumns: []*schema.Column{UsersColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// GroupOwnerColumns holds the columns for the "group_owner" table.
	GroupOwnerColumns = []*schema.Column{
		{Name: "group_id", Type: field.TypeUUID},
		{Name: "user_id", Type: field.TypeUUID},
	}
	// GroupOwnerTable holds the schema information for the "group_owner" table.
	GroupOwnerTable = &schema.Table{
		Name:       "group_owner",
		Columns:    GroupOwnerColumns,
		PrimaryKey: []*schema.Column{GroupOwnerColumns[0], GroupOwnerColumns[1]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "group_owner_group_id",
				Columns:    []*schema.Column{GroupOwnerColumns[0]},
				RefColumns: []*schema.Column{GroupsColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "group_owner_user_id",
				Columns:    []*schema.Column{GroupOwnerColumns[1]},
				RefColumns: []*schema.Column{UsersColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		CommentsTable,
		FilesTable,
		GroupsTable,
		GroupBudgetsTable,
		RequestsTable,
		RequestStatusTable,
		RequestTargetsTable,
		TagsTable,
		TransactionsTable,
		TransactionDetailsTable,
		UsersTable,
		GroupUserTable,
		GroupOwnerTable,
	}
)

func init() {
	CommentsTable.ForeignKeys[0].RefTable = RequestsTable
	FilesTable.ForeignKeys[0].RefTable = RequestsTable
	GroupBudgetsTable.ForeignKeys[0].RefTable = GroupsTable
	RequestsTable.ForeignKeys[0].RefTable = GroupsTable
	RequestStatusTable.ForeignKeys[0].RefTable = RequestsTable
	RequestTargetsTable.ForeignKeys[0].RefTable = RequestsTable
	TagsTable.ForeignKeys[0].RefTable = RequestsTable
	TagsTable.ForeignKeys[1].RefTable = TransactionsTable
	TransactionsTable.ForeignKeys[0].RefTable = GroupBudgetsTable
	TransactionsTable.ForeignKeys[1].RefTable = RequestsTable
	TransactionDetailsTable.ForeignKeys[0].RefTable = TransactionsTable
	UsersTable.ForeignKeys[0].RefTable = CommentsTable
	UsersTable.ForeignKeys[1].RefTable = RequestsTable
	UsersTable.ForeignKeys[2].RefTable = RequestStatusTable
	GroupUserTable.ForeignKeys[0].RefTable = GroupsTable
	GroupUserTable.ForeignKeys[1].RefTable = UsersTable
	GroupOwnerTable.ForeignKeys[0].RefTable = GroupsTable
	GroupOwnerTable.ForeignKeys[1].RefTable = UsersTable
}
