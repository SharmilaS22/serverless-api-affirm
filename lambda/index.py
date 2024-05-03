# lambda handler print hello world
import json
import boto3

def get_aff_handler(event, context):

    dynamodb = boto3.resource('dynamodb')
    table = dynamodb.Table('affirmations')
    response = table.scan()
    
    return {
        'statusCode': 200,
        "headers": {
            "Content-Type": "application/json"
        },
        'body': json.dumps(response['Items'])
    }


def post_aff_handler(event, context):

    dynamodb = boto3.resource('dynamodb')
    table = dynamodb.Table('affirmations')

    table.put_item(
        Item={
            'id': event['body']['id'],
            'affirmation': event['body']['affirmation']
        }
    )

    return {
        'statusCode': 201,
        "headers": {
            "Content-Type": "application/json"
        },
        'body': json.dumps({
            'message': 'Entry Added successfully!',
            'id': event['body']['id'],
        })
    }