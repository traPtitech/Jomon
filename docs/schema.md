# DB schema

**Jomon**の DB です。外部キー制約は全て`ON UPDATE reference_option`,`ON DELETE reference_option`共にデフォルト(`RESTRICT`)

## administrators

jomon の admin (会計の人：申請書更新等の権限)（admin のログはとりません）(権限剥奪の場合はレコードを削除します）

| Field   | Type        | Null | Key | Default | Extra | 説明など |
| ------- | ----------- | ---- | --- | ------- | ----- | -------- |
| trap_id | varchar(32) | NO   | PRI | NULL    |       |          |

## transactions

#### トランザクション

| Field      | Type     | Null | Key | Default           | Extra | 説明など                         |
| ---------- | -------- | ---- | --- | ----------------- | ----- | -------------------------------- |
| id         | char(36) | NO   | PRI | NULL              |       | uuid                             |
| created_at | datetime | NO   |     | CURRENT_TIMESTAMP | index | トランザクションが作成された時間 |

## transaction_details

#### トランザクションの詳細

| Field          | Type        | Null | Key | Default           | Extra | 説明など                                                                    |
| -------------- | ----------- | ---- | --- | ----------------- | ----- | --------------------------------------------------------------------------- |
| id             | char(36)    | NO   | PRI | NULL              |       | uuid                                                                        |
| transaction_id | char(36)    | NO   | MUL | NULL              | index |                                                                             |
| amount         | int(11)     | NO   |     | NULL              |       | 申請金額                                                                    |
| target         | varchar(64) | NO   |     | NULL              |       | 入金元 or 出金先(amount の正負で判定)                                       |
| request_id     | char(36)    | YES  | MUL | NULL              | index | 依頼への参照(NULL のときは依頼なし)**Parents:request.id**                   |
| group_id       | char(36)    | YES  | MUL | NULL              | index | グループへの参照(NULL のときはグループに所属していない)**Parents:group.id** |
| created_at     | datetime    | NO   |     | CURRENT_TIMESTAMP | index | トランザクションが作成/修正された時間                                       |

## transaction_tags

#### トランザクションのタグ

| Field          | Type     | Null | Key | Default           | Extra | 説明など                                           |
| -------------- | -------- | ---- | --- | ----------------- | ----- | -------------------------------------------------- |
| id             | char(36) | NO   | PRI | NULL              |       | 状態 ID uuid                                       |
| transaction_id | char(36) | NO   | MUL | NULL              | index | トランザクションへの参照**Parents:transaction.id** |
| tag_id         | char(36) | NO   | MUL | NULL              | index | タグへの参照**Parents:tag.id**                     |
| created_at     | datetime | NO   |     | CURRENT_TIMESTAMP |       | タグが追加された日時                               |

## requests

#### 依頼

新規、変更ごとに新しいレコードを作成。依頼の削除はできず、一度作ったら状態で管理

| Field      | Type        | Null | Key | Default           | Extra | 説明など             |
| ---------- | ----------- | ---- | --- | ----------------- | ----- | -------------------- |
| id         | varchar(36) | NO   | PRI | NULL              |       | uuid                 |
| created_by | varchar(32) | NO   |     | NULL              |       | traP ID              |
| amount     | int(11)     | NO   |     | NULL              |       | 申請金額             |
| created_at | datetime    | NO   |     | CURRENT_TIMESTAMP |       | 依頼が作成された時間 |

## request_statuses

#### 依頼の状態

状態の変更があるたびにレコードを作成。対応する依頼のレコード全ての`target`に対して`request_target`の paid_at が挿入されていたら`fully_repaid`に変更。新規の依頼ごとに新しいレコードを作成。request が作られた段階で作られる。reason は status を「submitted から fix_required」「submitted から rejected」「accepted から submitted」にするときに必要。作成者は「fix_required から submitted」にでき、admin は「submitted から rejected」「submitted から required」「fix_required から submitted」「submitted から accepted」「accepted から submitted(ただし、すでに払う/払われている人がいた場合、この操作は不可)」の操作が可能。

| Field      | Type        | Null | Key | Default           | Extra          | 説明など                           |
| ---------- | ----------- | ---- | --- | ----------------- | -------------- | ---------------------------------- |
| id         | int(11)     | NO   | PRI | NULL              | auto_increment | 状態 ID                            |
| request_id | varchar(36) | NO   | MUL | NULL              | index          | 依頼への参照**Parents:request.id** |
| created_by | varchar(32) | NO   |     | NULL              |                | 状態を変えた人の traPid            |
| status     | enum        | NO   |     | NULL              |                |                                    |
| reason     | text        | NO   |     | NULL              |                |                                    |
| created_at | datetime    | NO   |     | CURRENT_TIMESTAMP |                | 状態が更新された日時               |

## request_targets

#### 依頼の target

| Field      | Type        | Null | Key | Default           | Extra          | 説明など                           |
| ---------- | ----------- | ---- | --- | ----------------- | -------------- | ---------------------------------- |
| id         | int(11)     | NO   | PRI | NULL              | auto_increment |                                    |
| request_id | char(36)    | NO   | MUL | NULL              |                | 依頼への参照**Parents:request.id** |
| target     | varchar(64) | NO   |     | NULL              |                | 入金元 or 出金先                   |
| paid_at    | date        | YES  |     | NULL              |                | 払う/払われた日                    |
| created_at | datetime    | NO   |     | CURRENT_TIMESTAMP |                | request_target が作成された日時    |

## request_tags

#### 依頼のタグ

| Field      | Type     | Null | Key | Default           | Extra | 説明など                           |
| ---------- | -------- | ---- | --- | ----------------- | ----- | ---------------------------------- |
| id         | char(36) | NO   | PRI | NULL              |       | 状態 ID uuid                       |
| request_id | char(36) | NO   | MUL | NULL              | index | 依頼への参照**Parents:request.id** |
| tag_id     | char(36) | NO   | MUL | NULL              | index | タグへの参照**Parents:tag.id**     |
| created_at | datetime | NO   |     | CURRENT_TIMESTAMP |       | タグが追加された日時               |

## files

#### 依頼 id に対応するファイル

| Field      | Type     | Null | Key | Default           | Extra | 説明など                           |
| ---------- | -------- | ---- | --- | ----------------- | ----- | ---------------------------------- |
| id         | char(36) | NO   | PRI | NULL              |       | uuid                               |
| request_id | char(36) | NO   | MUL | NULL              |       | 依頼への参照**Parents:request.id** |
| mime_type  | text     | NO   |     | NULL              |       | フォーマット                       |
| created_at | datetime | NO   |     | CURRENT_TIMESTAMP |       | 登録された日時                     |
| deleted_at | datetime | YES  |     | NULL              |       | 削除された日時                     |

## comments

#### 依頼ごとのコメント

| Field      | Type        | Null | Key | Default           | Extra                       | 説明など                           |
| ---------- | ----------- | ---- | --- | ----------------- | --------------------------- | ---------------------------------- |
| id         | int(11)     | NO   | PRI | NULL              | auto_increment              | コメント ID                        |
| request_id | varchar(36) | NO   | MUL | NULL              | index                       | 依頼への参照**Parents:request.id** |
| created_by | varchar(32) | NO   |     | NULL              |                             |                                    |
| comment    | text        | NO   |     | NULL              |                             | コメント内容                       |
| created_at | datetime    | NO   |     | CURRENT_TIMESTAMP |                             | コメントが作成された日時           |
| updated_at | datetime    | NO   |     | CURRENT_TIMESTAMP | on update CURRENT_TIMESTAMP | コメントが更新された日時           |
| deleted_at | datetime    | YES  |     | NULL              |                             | コメントが削除された日時           |

## groups

#### グループ

| Field       | Type        | Null | Key | Default           | Extra                       | 説明など                                                           |
| ----------- | ----------- | ---- | --- | ----------------- | --------------------------- | ------------------------------------------------------------------ |
| id          | char(36)    | NO   | PRI | NULL              |                             | uuid                                                               |
| created_at  | datetime    | NO   |     | CURRENT_TIMESTAMP |                             | 登録された日時                                                     |
| updated_at  | datetime    | NO   |     | CURRENT_TIMESTAMP | on update CURRENT_TIMESTAMP | 変更された日時                                                     |
| deleted_at  | datetime    | YES  |     | NULL              |                             | 削除された日時                                                     |
| name        | varchar(64) | NO   |     | NULL              |                             | グループ名                                                         |
| description | text        | NO   |     | NULL              |                             | グループの説明                                                     |
| budget      | int(11)     | YES  |     | NULL              |                             | 予算額 (あえて非正規化してる) 履歴は group_budget テーブルを参照。 |

## group_budgets

#### グループの予算

| Field      | Type     | Null | Key | Default           | Extra | 説明など       |
| ---------- | -------- | ---- | --- | ----------------- | ----- | -------------- |
| id         | char(36) | NO   | PRI | NULL              |       | uuid           |
| created_at | datetime | NO   |     | CURRENT_TIMESTAMP |       | 登録された日時 |
| group_id   | char(36) | NO   | MUL | NULL              | index | uuid           |
| amount     | int(11)  | NO   |     | NULL              |       | 予算額         |

## group_users

#### グループのユーザー

| Field      | Type        | Null | Key | Default           | Extra | 説明など       |
| ---------- | ----------- | ---- | --- | ----------------- | ----- | -------------- |
| id         | char(36)    | NO   | PRI | NULL              |       | uuid           |
| created_at | datetime    | NO   |     | CURRENT_TIMESTAMP |       | 登録された日時 |
| group_id   | char(36)    | NO   | MUL | NULL              | index | uuid           |
| user_id    | varchar(32) | NO   |     | NULL              |       | traPID         |

## group_owners

#### グループのオーナー

| Field      | Type        | Null | Key | Default           | Extra | 説明など       |
| ---------- | ----------- | ---- | --- | ----------------- | ----- | -------------- |
| id         | char(36)    | NO   | PRI | NULL              |       | uuid           |
| created_at | datetime    | NO   |     | CURRENT_TIMESTAMP |       | 登録された日時 |
| group_id   | char(36)    | NO   | MUL | NULL              | index | uuid           |
| owner      | varchar(32) | NO   | MUL | NULL              |       | traPID         |

## tags

#### タグ

| Field       | Type        | Null | Key | Default           | Extra                       | 説明など       |
| ----------- | ----------- | ---- | --- | ----------------- | --------------------------- | -------------- |
| id          | char(36)    | NO   | PRI | NULL              |                             | uuid           |
| created_at  | datetime    | NO   |     | CURRENT_TIMESTAMP |                             | 登録された日時 |
| updated_at  | datetime    | NO   |     | CURRENT_TIMESTAMP | on update CURRENT_TIMESTAMP | 変更された日時 |
| deleted_at  | datetime    | YES  |     | NULL              |                             | 削除された日時 |
| name        | varchar(64) | NO   |     | NULL              |                             | タグ名         |
| description | text        | NO   |     | NULL              |                             | タグの説明     |
