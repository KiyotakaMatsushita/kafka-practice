# kafka-practice

このリポジトリは、Kafka と Kafka クライアントを利用した動作検証のためのプロジェクトです。

## Kafka Broker の起動方法

1. 以下のコマンドを実行して Kafka Broker を起動します：

   ```bash
   docker-compose up -d
   ```

2. Broker が正常に起動したことを確認します：

   ```bash
   docker-compose ps
   ```

## Kafka Topic の作成方法

1. Kafka Broker のコンテナに接続します：

   ```bash
   docker exec --workdir /opt/kafka/bin/ -it broker sh
   ```

2. 以下のコマンドで Topic を作成します：

   ```bash
   ./kafka-topics.sh --create --topic default_topic --bootstrap-server localhost:9092 --partitions 1 --replication-factor 1
   ```

## Kafka Producer から topic メッセージの作成方法

### 正常なメッセージの作成方法

1. Kafka Broker のコンテナに接続します：

   ```bash
   docker exec --workdir /opt/kafka/bin/ -it broker sh
   ```

2. 以下のコマンドで Producer を起動し、メッセージを送信します：

   ```bash
   ./kafka-console-producer.sh --bootstrap-server localhost:9092 --topic default_topic
   ```

3. プロンプトが表示されたら、メッセージを入力し、Enter キーを押して送信します。

### メッセージの作成方法

1. 上記の手順で Producer を起動した後、以下の JSON メッセージを入力してパニックを発生させます：

   ```bash
   > panic-message
   ```

2. 上記の手順で Producer を起動した後、以下の JSON の正常なメッセージを入力します：

   ```bash
   > {"hoge": "fuga"}
   ```

## Kafka Consumer Group の message 確認方法

1. Kafka Broker のコンテナに接続します：

   ```bash
   docker exec --workdir /opt/kafka/bin/ -it broker sh
   ```

2. Consumer Group を起動してメッセージを確認します：

   ```bash
   ./kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic default_topic --from-beginning
   ```

3. Consumer Group の詳細情報を確認するには、以下のコマンドを使用します：

   ```bash
   ./kafka-consumer-groups.sh --bootstrap-server localhost:9092 --describe --group default_group
   ```

   このコマンドは、指定された Consumer Group（この場合は`default_group`）の詳細情報を表示します。これには、グループのメンバー、消費しているトピックとパーティション、現在のオフセット、ラグなどの情報が含まれます。

## Docker Go-app の起動方法

1. 以下のコマンドで Go-app を起動します：

   ```bash
   docker-compose up go-app
   ```

2. アプリケーションのログを確認します：

   ```bash
   docker-compose logs -f go-app
   ```
