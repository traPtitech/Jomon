# DB schema

**Jomon**のDBです。外部キー制約は **（※）** を除いて全て`ON UPDATE reference_option`,`ON DELETE reference_option`共にデフォルト(`RESTRICT`)です。申請書リスト取得時の処理を高めるために、appolicationsに'applications_details_id','states_logs_id'を追加したことで相互参照が起こります。よって以上二つ(※）については内容としてはMULですが、DMSによっては制限しないことにします。

## administrators

jomonのadmin (会計の人：申請書更新等の権限)（adminのログはとりません）(権限剥奪の場合はレコードを削除します）

| Field            | Type     | Null | Key | Default | Extra | 説明など |
| ---------------- | -------- | ---- | --- | ------- | ----- | -------- |
| admin_trap_id     | varchar(32) | NO   | PRI | _NULL_  |

## applications

同一の経費精算書類の情報を持ちます。削除はできません。

| Field            | Type       | Null | Key | Default           | Extra          | 説明など                                                                                                       |
| ---------------- | ---------- | ---- | --- | ----------------- | -------------- | -------------------------------------------------------------------------------------------------------------- |
| id          | char(36) | NO   | PRI | _NULL_  |  |uuid|
| applications_details_id          | int(11) | NO   | MUL | _NULL_  || 経費申請詳細の最新id**Parents:applications_details.id** **（※）** |
| states_logs_id          | int(11) | NO   | MUL | _NULL_  || 状態の最新id**Parents:states_logs.id**　**（※）**  |
| create_user_trap_id      | varchar(32) | NO   | MUL | _NULL_  |           | 申請者のtraPid |
| created_at       | timestamp  | NO   |     | CURRENT_TIMESTAMP |       | 申請書が作成された日時 |

## applications_details

経費精算申請（新規、変更ごとに新しいレコードが作られます。申請の削除はできず、一度作ったら必ずいずれかのstateに当てはまります。)

| Field            | Type       | Null | Key | Default           | Extra          | 説明など                                                                                                       |
| ---------------- | ---------- | ---- | --- | ----------------- | -------------- | -------------------------------------------------------------------------------------------------------------- |
| id          | int(11) | NO   | PRI | _NULL_  | auto_increment |
|application_id|char(36)|NO|MUL|_NULL_||経費精算申請ごとにつくid **parents:applications.id**|
| update_user_trap_id      | varchar(32) | NO   | MUL | _NULL_  |           | 変更者（初めは申請者）のtraPid |
| type             | tinyint(4)   | NO   |     | _NULL_            |                | どのタイプの申請か (0(Club), 1(Contest), 2(Event), 3(Public)) |
| title        | text      | NO  |     | _NULL_||        申請の目的、概要(大会名など) |
| remarks       | text      | YES  |     | _NULL_ |           |   備考（購入したものの概要、旅程、乗車区間など） |
| amount | int(11)    | NO  |     | _NULL_    |         |申請金額    |
| bought_at       | timestamp  | NO   |     |  |       | お金を使った日  |
| created_at       | timestamp  | NO   |     | CURRENT_TIMESTAMP |       | 申請書が作成（変更）された日時  |

## return_users

申請idにつき、誰に返金されるか　(払い戻し対象者の変更ログは残りません)(現在usertableがないためtraPidはtraQ(できればportal)のapiをたたきます。)(変更時には対応する`application_id`のレコードすべてを削除して、新しいレコードを追加します。)

| Field            | Type       | Null | Key | Default           | Extra          | 説明など                                                                                                       |
| ---------------- | ---------- | ---- | --- | ----------------- | -------------- | -------------------------------------------------------------------------------------------------------------- |
| id          | int(11) | NO   | PRI | _NULL_  |auto_increment|  |
| application_id          | char(36) | NO   | MUL | _NULL_  || 申請書のid |
| returned_to_user_trap_id      | varchar(32) | NO   | MUL | _NULL_  |           | 払い戻される人のtraPid |
| returned_by_user_trap_id      | varchar(32) | YES   | MUL | _NULL_  |           | お金を渡した人のtraPid |
| returned_at          | timestamp | YES   |  | _NULL_  | |払い戻された日  |

## applications_images

申請idにつき、対応する画像　(画像変更ログは残りません。)(変更時には対応する`application_id`のレコードすべてを削除して、新しいレコードを追加します。)

| Field            | Type       | Null | Key | Default           | Extra          | 説明など                                                                                                       |
| ---------------- | ---------- | ---- | --- | ----------------- | -------------- | -------------------------------------------------------------------------------------------------------------- |
| id          |char(36) | NO   | PRI | _NULL_  || uuid |
| application_id          | char(36) | NO   | MULL | _NULL_  || 申請書のid |

## states_logs

状態の記録（状態の変更があるたびにレコードを追加します。）(初めて申請書が作られたときも0をレコードとして入れます。）（理由の変更、削除はできません。)(stateの`3`は`return_users`に依存していて、全員が`true`となった時に変えてください。)

| Field            | Type       | Null | Key | Default           | Extra          | 説明など                                                                                                       |
| ---------------- | ---------- | ---- | --- | ----------------- | -------------- | -------------------------------------------------------------------------------------------------------------- |
| id          | int(11) | NO   | PRI | _NULL_  |auto_increment|  |
| application_id          | char(36) | NO   | MUL | _NULL_  || 申請書のid **parents:applications.id**|
| update_user_trap_id      | varchar(32) | NO   |  | _NULL_  |           | 状態を変えた人のtraPid |
| to_state     | tinyint(4) | NO   |     | 0                 |                | どの状態へ変えたか (0(submitted) ,1(fix_required), 2(accepted), 3(fully_returned), 4(rejected))                                                                                 |
| reason     |text | YES  |     | _NULL_                 |                | 状態を変えたとき状態の変え方によってコメントをつけられたり付けられなかったりします。（swagger参照) |
| created_at       | timestamp  | NO   |     | CURRENT_TIMESTAMP |                | 状態が更新された日時                                                                                                  |

## comments

申請書ごとのコメント（コメントの変更、削除は対応するレコードを変更することで行います。そのため変更前の状態履歴は残りません。）

| Field            | Type      | Null | Key | Default           | Extra          | 説明など                                            |
| ---------------- | --------- | ---- | --- | ----------------- | -------------- | --------------------------------------------------- |
| id      | int(11)   | NO   | PRI | _NULL_            | auto_increment | コメントIＤ |
| application_id | char(36)  | NO   | MUL | _NULL_            |                | どの申請書へのコメントか **Parents:applications.id**                          |
| user_trap_id      | varchar(32)  | NO  | MUL | _NULL_            |                | コメントした人の traPID                                     |
| comment       |  text    | NO  |     | _NULL_            |       |コメント内容そのもの                                       |
| created_at     | timestamp | NO   |     | CURRENT_TIMESTAMP |                | コメントが作成された日時                                                                                              |
| updated_at     | timestamp |  NO  |     | CURRENT_TIMESTAMP |    on update CURRENT_TIMESTAMP            | コメントが更新された日時                                                                                              |
| deleted_at     | timestamp |  YES  |     | NULL |                | コメントが削除された日時                                                                                              |
