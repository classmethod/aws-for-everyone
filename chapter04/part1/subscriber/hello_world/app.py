import base64
import boto3
import json
import logging
import os
from botocore.exceptions import ClientError
from datetime import datetime, timezone, timedelta

logger = logging.getLogger()
logger.setLevel(logging.INFO)

dynamodb = boto3.resource('dynamodb')


def lambda_handler(event, context):
    records_length = len(event['records'])
    logger.info(f'event: {json.dumps(event)}')
    logger.info(f'records_length: {records_length}')

    transformed_data = []

    for index, record in enumerate(event['records']):
        data = {}
        result = 'Ok'
        try:
            # 現在処理しているデータ（ログ用）
            log_header = f'[{index + 1}/{records_length}]'

            # 1件分のデータを取得する
            payload = json.loads(base64.b64decode(record['data']))
            logger.info(f'{log_header} payload: {json.dumps(payload)}')

            # DynamoDBからETCゲートの情報を取得する
            device_item = get_item(payload['serialNumber'])

            # データを変換（追加）する
            payload['feeStationNumber'] = device_item['feeStationNumber']
            payload['feeStationName'] = device_item['feeStationName']
            payload['gateNumber'] = int(device_item['gateNumber'])
            payload['timestring'] = convert_iso_format(payload['timestamp'])
        except ClientError as e:
            error_message = e.response['Error']['Message']
            logger.error(f'{log_header} DynamoDB ClientError: {error_message}')
            result = 'Ng'
        except Exception as e:
            logger.error(f'{log_header} Transform failed: {e}')
            result = 'Ng'

        # Firehoseに戻すデータを作る
        data = json.dumps(payload) + '\n'
        logger.info(f'{log_header} transformed: {data}')

        data_utf8 = data.encode('utf-8')
        transformed_data.append({
            'recordId': record['recordId'],
            'result': result,
            'data': base64.b64encode(data_utf8).decode('utf-8')
        })

    logger.info('finish transform.')

    return {
        'records': transformed_data
    }


def get_item(serial_number:int) -> dict:
    table_name = os.getenv('ETC_GATE_MANAGEMENT_TABLE_NAME')
    table = dynamodb.Table(table_name)

    # DynamoDBからETCゲートの情報を取得する
    res = table.get_item(Key={
        'serialNumber': serial_number
    })

    return res['Item']


def convert_iso_format(timestamp:int) -> str:
    # JSTとするので+9時間する
    tz = timezone(timedelta(hours=9))
    timestamp =  datetime.fromtimestamp(timestamp / 1000, tz)
    return timestamp.isoformat(timespec='milliseconds')
