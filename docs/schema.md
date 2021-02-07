# DB schema

**Jomon**のDBです。外部キー制約は全て`ON UPDATE reference_option`,`ON DELETE reference_option`共にデフォルト(`RESTRICT`)

## administrators

jomonのadmin (会計の人：申請書更新等の権限)（adminのログはとりません）(権限剥奪の場合はレコードを削除します）

| Field            | Type     | Null | Key | Default | Extra | 説明など |
| ---------------- | -------- | ---- | --- | ------- | ----- | -------- |
| trap_id     | varchar(32) | NO   | PRI | _NULL_  |

## request
#### 依頼
新規、変更ごとに新しいレコードを作成。依頼の削除はできず、一度作ったら状態で管理

| Field      | Type        | Null | Key | Default | Extra | 説明など                           |
| ---------- | ----------- | ---- | --- | ------- | ----- | ---------------------------------- |
| id         | varchar(36) | NO   | PRI | NULL    |       | uuid                               |
| created_by | varchar(32) | NO   |     | NULL    |       | traP ID                            |
| amount     | int(11)     | NO   |     | NULL    |       | 申請金額                           |
|     created_at       |    datetime         |   NO   |     |     CURRENT_TIMESTAMP    |       |               依頼が作成された時間                     |

## request_status
#### 依頼の状態
状態の変更があるたびにレコードを作成。対応する依頼のレコード全ての`target`に対して`request_target`のpaid_atが挿入されていたら`fully_repaid`に変更

| Field      | Type        | Null | Key | Default           | Extra          | 説明など                           |
| ---------- | ----------- | ---- | --- | ----------------- | -------------- | ---------------------------------- |
| id         | int(11)     | NO   | PRI | NULL              | auto_increment | 状態ID                         |
| request_id | varchar(36) | NO   | MUL | NULL              | index          | 依頼への参照**Parents:request.id** |
| created_by      | varchar(32) | NO   |  | NULL  |           | 状態を変えた人のtraPid |
| status     | enum        | NO   |     | NULL              |                |                                    |
| reason     |text | NO  |     | NULL                 |                |  |
| created_at | datetime    | NO   |     | CURRENT_TIMESTAMP |                | 状態が更新された日時            |


## transaction
#### 入出金
実際にすでに行われた入出金をすべて記録。新規ごとに新しいレコードを作成。requestが作られた段階で作られる。reasonはstatusを「submittedからfix_required」「submittedからrejected」「acceptedからsubmitted」にするときに必要。作成者は「fix_requiredからsubmitted」にでき、adminは「submittedからrejected」「submittedからrequired」「fix_requiredからsubmitted」「submittedからaccepted」「acceptedからsubmitted(ただし、すでに払う/払われている人がいた場合、この操作は不可)」の操作が可能。

| Field      | Type        | Null | Key   | Default | Extra | 説明など                           |
| ---------- | ----------- | ---- | ----- | ------- | ----- | ---------------------------------- |
| id         | char(36)    | NO   | PRI   | NULL    |       | uuid                               |
| amount     | int(11)     | NO   |       | NULL    |       | 申請金額                           |
| target     | varchar(64) | NO   |       | NULL    |       | 入金元or出金先(amountの正負で判定) |
| request_id | varchar(36) | YES  | MUL | NULL    |   index    | 依頼への参照(NULLのときは依頼なし)**Parents:request.id** |
| created_at           | datetime            |  NO    |       |   CURRENT_TIMESTAMP      | index      |                           トランザクションが作成された時間         |


## request_target
#### 依頼のtarget

| Field                  | Type        | Null | Key | Default           | Extra          | 説明など                    |
| ---------------------- | ----------- | ---- | --- | ----------------- | -------------- | --------------------------- |
| id                     | int(11)     | NO   | PRI | NULL            | auto_increment |                             |
| request_id         | char(36)    | NO   | MUL | NULL            |                | 依頼への参照**Parents:request.id**                  |
| target | varchar(64) | NO   |  | NULL            |                | 入金元or出金先      |
| paid_at              | date        | YES  |     | NULL            |                | 払う/払われた日              |
| created_at             | datetime    | NO   |     | CURRENT_TIMESTAMP |                | request_targetが作成された日時 |


## file
#### 依頼idに対応するファイル

| Field      | Type     | Null | Key  | Default           | Extra | 説明など                           |
| ---------- | -------- | ---- | ---- | ----------------- | ----- | ---------------------------------- |
| id         | char(36) | NO   | PRI  | NULL              |       | uuid                               |
| request_id | char(36) | NO   | MUL | NULL              |       | 依頼への参照**Parents:request.id** |
| mime_type  | text     | NO   |      | NULL              |       | フォーマット                       |
| created_at | datetime | NO   |      | CURRENT_TIMESTAMP |       | 登録された日時                     |
| deleted_at | datetime    | YES  |     | NULL              |                             | 削除された日時           |


## comment
#### 依頼ごとのコメント

| Field      | Type        | Null | Key | Default           | Extra                       | 説明など                           |
| ---------- | ----------- | ---- | --- | ----------------- | --------------------------- | ---------------------------------- |
| id         | int(11)     | NO   | PRI | NULL              | auto_increment              | コメントID                         |
| request_id | varchar(36) | NO   | MUL | NULL              | index                       | 依頼への参照**Parents:request.id** |
| created_by | varchar(32) | NO   |     | NULL              |                             |                                    |
| comment    | text        | NO   |     | NULL              |                             | コメント内容                       |
| created_at | datetime    | NO   |     | CURRENT_TIMESTAMP |                             | コメントが作成された日時           |
| updated_at | datetime    | NO   |     | CURRENT_TIMESTAMP | on update CURRENT_TIMESTAMP | コメントが更新された日時           |
| deleted_at | datetime    | YES  |     | NULL              |                             | コメントが削除された日時           |
