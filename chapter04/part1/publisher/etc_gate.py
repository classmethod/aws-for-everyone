import json
import os
import random
import sys
import time
from AWSIoTPythonSDK.MQTTLib import AWSIoTMQTTClient

IOT_CORE_ENDPOINT = 'xxxxx-ats.iot.ap-northeast-1.amazonaws.com'
PORT = 8883
TOPIC_NAME = 'etc_gate/passing/car'
QOS = 0

ROOT_CA_FILE = './AmazonRootCA1.pem'

ETC_GATE_INFO = {
    '1111ABCD': {
        'client_id': 'etc_gate_1111ABCD',
        'certificate_file': './etc_gate_1111ABCD_certificate.pem',
        'private_key_file': './etc_gate_1111ABCD_certificate.private',
        'rate': {
            'A': 97.0,
            'B': 0.1,
            'C': 0.1,
            'D': 2.8
        }
    },
    '2222EFGH': {
        'client_id': 'etc_gate_2222EFGH',
        'certificate_file': './etc_gate_2222EFGH_certificate.pem',
        'private_key_file': './etc_gate_2222EFGH_certificate.private',
        'rate': {
            'A': 97.0,
            'B': 0.1,
            'C': 0.1,
            'D': 2.8
        }
    },
    '3333IJKL': {
        'client_id': 'etc_gate_3333IJKL',
        'certificate_file': './etc_gate_3333IJKL_certificate.pem',
        'private_key_file': './etc_gate_3333IJKL_certificate.private',
        'rate': {
            'A': 90.0,
            'B': 4.0,
            'C': 4.0,
            'D': 2.0
        }
    },
}

FINISH_FILE = './finish.txt'

def main(serial_number: str) -> None:
    init()

    client_id = ETC_GATE_INFO[serial_number]['client_id']
    certificate_file = ETC_GATE_INFO[serial_number]['certificate_file']
    private_key_file = ETC_GATE_INFO[serial_number]['private_key_file']

    # IoT Coreに接続する
    # https://github.com/aws/aws-iot-device-sdk-python
    # https://s3.amazonaws.com/aws-iot-device-sdk-python-docs/sphinx/html/index.html
    client = AWSIoTMQTTClient(client_id)
    client.configureEndpoint(IOT_CORE_ENDPOINT, PORT)
    client.configureCredentials(
        ROOT_CA_FILE,
        private_key_file,
        certificate_file)
    client.connect()

    while True:
        data = create_data(serial_number)

        # IoT CoreのトピックにPublishする
        client.publish(TOPIC_NAME, json.dumps(data), QOS)

        time.sleep(1)
        if is_finish():
            break

def init() -> None:
    # もし finish.txt があるなら削除しておく
    if os.path.isfile(FINISH_FILE):
        os.remove(FINISH_FILE)

def create_data(serial_number: str) -> dict:
    num = random.uniform(0, 100)
    rate = ETC_GATE_INFO[serial_number]['rate']

    (rate_border_a, rate_border_b, rate_border_c) = get_rate_border(rate)

    param = {}
    if 0 <= num and num < rate_border_a:
        # パターンA
        param['open'] = True
        param['payment'] = True
    elif rate_border_a <= num < rate_border_b:
        # パターンB
        param['open'] = True
        param['payment'] = False
    elif rate_border_b <= num < rate_border_c:
        # パターンC
        param['open'] = False
        param['payment'] = True
    else:
        # パターンD
        param['open'] = False
        param['payment'] = False

    return {
        'serialNumber': '1111ABCD',
        'timestamp': int(time.time() * 1000),
        'open': param['open'],
        'payment': param['payment']
    }

def get_rate_border(rate: dict) -> (float, float, float):
    return (
        rate['A'],
        rate['A'] + rate['B'],
        rate['A'] + rate['B'] + rate['C']
    )

def is_finish():
    if os.path.isfile(FINISH_FILE):
        return True
    return False


if __name__ == '__main__':
    args = sys.argv
    if len(args) == 2:
        main(args[1])
