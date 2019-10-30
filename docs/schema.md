# DB schema

### administrators

jomonのadmin (会計の人：申請書更新等の権限)（adminのログはとりません）

| Field            | Type     | Null | Key | Default | Extra | 説明など |
| ---------------- | -------- | ---- | --- | ------- | ----- | -------- |
| admin_traqid     | char(30) | NO   | PRI | _NULL_  |


### applications

経費精算申請（新規、変更ごとに新しレコードが作られます。申請書の削除はできず、一度作ったら必ずいずれかのstateに当てはまります。created_atが新しい順にかつapplications_idが一つとなるようにすれば最新の状態が得られます。）(画像は複数可とするか？可とするなら別個テーブルを作ります。)

| Field            | Type       | Null | Key | Default           | Extra          | 説明など                                                                                                       |
| ---------------- | ---------- | ---- | --- | ----------------- | -------------- | -------------------------------------------------------------------------------------------------------------- |
| id          | int(11) | NO   | PRI | _NULL_  | auto_increment |
|applications_id|int(11)|NO||_NULL_||経費精算申請ごとにつくid|
| create_user_traq_id      | char(30) | NO   | MUL | _NULL_  |           | 申請者のtraQid |
| type             | char(20)   | NO   |     | _NULL_            |                | どのタイプの申請か ("Club", "Contest", "Event", "Public") |
| title        | text      | NO  |     | _NULL_||        申請の目的、概要(大会名など) |
| remarks       | text      | YES  |     | _NULL_ |           |   備考（購入したものの概要、旅程、乗車区間など） |
| image_name | text | YES   |     |_NULL_       |       | 領収書等の画像   |
| ammount | int(11)    | NO  |     | _NULL_    |         |申請金額    |                     
| created_at       | timestamp  | NO   |     | CURRENT_TIMESTAMP |       | 申請書が作成された日時      |



### return_users

申請idにつき、誰に返金されるか　(払い戻し対象者の変更ログは残りません)(現在usertableがないためtraQidはtraQのapiをたたく必要がありそう。)

| Field            | Type       | Null | Key | Default           | Extra          | 説明など                                                                                                       |
| ---------------- | ---------- | ---- | --- | ----------------- | -------------- | -------------------------------------------------------------------------------------------------------------- |
| application_id          | int(11) | NO   | PRI | _NULL_  |auto_increment| 申請書のid |
| reimbursed_user_traq_id      | char(30) | NO   | MUL | _NULL_  |           | 払い戻さる人のtraQid |

### state_logs

状態の記録（状態の変更があるたびにレコードを追加）

| Field            | Type       | Null | Key | Default           | Extra          | 説明など                                                                                                       |
| ---------------- | ---------- | ---- | --- | ----------------- | -------------- | -------------------------------------------------------------------------------------------------------------- |
| id          | int(11) | NO   | PRI | _NULL_  |auto_increment|  |
| application_id          | int(11) | NO   | MUL | _NULL_  || 申請書のid |
| change_user_traq_id      | char(30) | NO   |  | _NULL_  |           | 状態を変えた人のtraQid |
| to_state     | tinyint(4) | NO   |     | 0                 |                | どの状態へ変えたか (申請済み(0) ,却下(1),要修正(2),許可済み(3),返金済み(4))                                                                                 |
| update_at       | timestamp  | NO   |     | CURRENT_TIMESTAMP |                | 状態が更新された日時                                                                                                  |



### comments

申請書ごとのコメント（コメントの変更、削除はできません）

| Field            | Type      | Null | Key | Default           | Extra          | 説明など                                            |
| ---------------- | --------- | ---- | --- | ----------------- | -------------- | --------------------------------------------------- |
| comment_id      | int(11)   | NO   | PRI | _NULL_            | auto_increment | コメントIＤ |
| application_id | int(11)   | NO   | MUL | _NULL_            |                | どの申請書へのコメントか                            |
| user_traq_id      | char(30)  | NO  | MUL | _NULL_            |                | コメントした人の traQID                                     |
| comment       |  char(36)    | NO  |  MUL   | _NULL_            |       | コメントのtraQでのid                                          |
| created_at     | timestamp | NO   |     | CURRENT_TIMESTAMP |                | コメントが作成された日時                                                                                              |

