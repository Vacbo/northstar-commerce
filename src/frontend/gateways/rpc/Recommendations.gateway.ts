// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

import { ChannelCredentials, Metadata } from '@grpc/grpc-js';
import { ListRecommendationsResponse, RecommendationServiceClient } from '../../protos/demo';

const { RECOMMENDATION_ADDR = '' } = process.env;

const client = new RecommendationServiceClient(RECOMMENDATION_ADDR, ChannelCredentials.createInsecure());

const DEADLINE_SECONDS = 5;

const RecommendationsGateway = () => ({
  listRecommendations(userId: string, productIds: string[]) {
    return new Promise<ListRecommendationsResponse>((resolve, reject) => {
      const deadline = new Date();
      deadline.setSeconds(deadline.getSeconds() + DEADLINE_SECONDS);
      const metadata = new Metadata();
      metadata.set('grpc-timeout', `${DEADLINE_SECONDS}S`);
      client.listRecommendations(
        { userId, productIds },
        metadata,
        { deadline },
        (error, response) =>
          error ? reject(error) : resolve(response)
      );
    });
  },
});

export default RecommendationsGateway();
