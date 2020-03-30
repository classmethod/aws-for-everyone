<?php

echo 'ECSタスク定義（コンテナ定義）の環境変数<br />';
echo 'app-alb-fargate.ymlのType: AWS::ECS::TaskDefinitionのEnvironment、およびSecretsで指定した内容が環境変数として取得されます。<br />';

echo '<p>DB関連環境変数一覧</p>';

$DBHOST = getenv('DBHOST');
$DB = getenv('DB');
$DBUSER = getenv('DBUSER');
$DBPASSWORD = getenv('DBPASSWORD');

echo 'DBHOST='.$DBHOST.'<br />';
echo 'DB='.$DB.'<br />';
echo 'DBUSER='.$DBUSER.'<br />';
echo 'DBPASSWORD='.$DBPASSWORD.'<br />';

try {
	// DB接続処理
	$dsn = 'mysql:dbname='.$DB.';host='.$DBHOST.';charset=utf8mb4';
	$pdo = new PDO($dsn, $DBUSER, $DBPASSWORD);
	
	$stmt = $pdo -> query('SELECT id, name FROM handsonUser');
	$stmt -> execute();

	//テーブル内容表示
	echo '<p>handsonUserテーブル内容</p>';
	while ($row = $stmt->fetch()) {
    printf("id:%s, name:%s<br />\n", $row['id'], $row['name']);
	}

} catch (PDOException $e)  {
	exit($e->getMessage()); 
}

?>
