// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

import { ChannelCredentials, credentials } from '@grpc/grpc-js';
import { ListRecommendationsResponse, RecommendationServiceClient } from '../../protos/demo';

const { RECOMMENDATION_ADDR = '' } = process.env;

const client = new RecommendationServiceClient(RECOMMENDATION_ADDR, ChannelCredentials.createInsecure());

const RecommendationsGateway = () => ({
  listRecommendations(userId: string, productIds: string[]) {
    return new Promise<ListRecommendationsResponse>((resolve, reject) => {
      const deadline = new Date();
      deadline.setSeconds(deadline.getSeconds() + 5);
      client.listRecommendations({ userId, productIds }, { credentials: credentials.createInsecure(), deadline }, (error, response) =>
        error ? reject(error) : resolve(response)
      );
    });
  },
});

export default RecommendationsGateway();
