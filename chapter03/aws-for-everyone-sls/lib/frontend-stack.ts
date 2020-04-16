import {Bucket} from '@aws-cdk/aws-s3';
import {CfnOutput, Construct, Duration, RemovalPolicy, Stack, StackProps} from "@aws-cdk/core";
import {BucketDeployment, Source} from '@aws-cdk/aws-s3-deployment';
import {CloudFrontWebDistribution, PriceClass, OriginAccessIdentity} from '@aws-cdk/aws-cloudfront'
import {CanonicalUserPrincipal, Effect, PolicyStatement} from '@aws-cdk/aws-iam';

export class FrontendStack extends Stack {
    constructor(scope: Construct, id: string, props?: StackProps) {
        super(scope, id, props);

        // Bucketの作成
        const websiteBucket = new Bucket(this, 'Website', {
            removalPolicy: RemovalPolicy.DESTROY
        });

        // OAIの作成
        const OAI = new OriginAccessIdentity(this, 'OAI');

        // Bucket Policy(OAIに関して)を作成
        const webSiteBucketPolicyStatement = new PolicyStatement({
            effect: Effect.ALLOW,
            actions: ['s3:GetObject'],
            resources: [`${websiteBucket.bucketArn}/*`],
            principals: [
                new CanonicalUserPrincipal(OAI.cloudFrontOriginAccessIdentityS3CanonicalUserId)
            ]
        });
        websiteBucket.addToResourcePolicy(webSiteBucketPolicyStatement);

        const distribution = new CloudFrontWebDistribution(this, 'WebsiteDistribution', {
            originConfigs: [
                {
                    s3OriginSource: {
                        s3BucketSource: websiteBucket,
                        originAccessIdentity: OAI
                    },
                    behaviors: [{
                        isDefaultBehavior: true,
                        minTtl: Duration.seconds(0),
                        maxTtl: Duration.seconds(0),
                        defaultTtl: Duration.seconds(0),
                    }]
                }
            ],
            errorConfigurations: [
                {
                    errorCode: 403,
                    responsePagePath: '/index.html',
                    responseCode: 200,
                    errorCachingMinTtl: 0,
                },
                {
                    errorCode: 404,
                    responsePagePath: '/index.html',
                    responseCode: 200,
                    errorCachingMinTtl: 0,
                }
            ],
            priceClass: PriceClass.PRICE_CLASS_200
        });

        new BucketDeployment(this, 'DeployWebsite', {
            sources: [Source.asset('src/frontend/dist')],
            destinationBucket: websiteBucket,
            distribution: distribution,
            distributionPaths: ['/*']
        });

        new CfnOutput(this, 'URL', {value: `https://${distribution.domainName}/`})
    }
}
