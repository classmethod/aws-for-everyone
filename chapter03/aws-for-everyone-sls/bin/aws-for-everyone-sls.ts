#!/usr/bin/env node
import 'source-map-support/register';
import {App} from '@aws-cdk/core';
import {BackendStack} from "../lib/backend-stack";
import {FrontendStack} from '../lib/frontend-stack';

const app = new App();
const region: string = 'ap-northeast-1';
new BackendStack(app, 'BackendStack', {env: {region: region}});
new FrontendStack(app, 'FrontendStack', {env:{region: region}});
