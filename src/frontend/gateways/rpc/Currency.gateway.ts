// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

import { ChannelCredentials, Metadata } from '@grpc/grpc-js';
import { GetSupportedCurrenciesResponse, CurrencyServiceClient, Money } from '../../protos/demo';

const { CURRENCY_ADDR = '' } = process.env;

const client = new CurrencyServiceClient(CURRENCY_ADDR, ChannelCredentials.createInsecure());

const DEADLINE_SECONDS = 5;

const CurrencyGateway = () => ({
  convert(from: Money, toCode: string) {
    return new Promise<Money>((resolve, reject) => {
      const deadline = new Date();
      deadline.setSeconds(deadline.getSeconds() + DEADLINE_SECONDS);
      const metadata = new Metadata();
      metadata.set('grpc-timeout', `${DEADLINE_SECONDS}S`);
      client.convert({ from, toCode }, metadata, { deadline }, (error, response) =>
        error ? reject(error) : resolve(response)
      );
    });
  },
  getSupportedCurrencies() {
    return new Promise<GetSupportedCurrenciesResponse>((resolve, reject) => {
      const deadline = new Date();
      deadline.setSeconds(deadline.getSeconds() + DEADLINE_SECONDS);
      const metadata = new Metadata();
      metadata.set('grpc-timeout', `${DEADLINE_SECONDS}S`);
      client.getSupportedCurrencies({}, metadata, { deadline }, (error, response) =>
        error ? reject(error) : resolve(response)
      );
    });
  },
});

export default CurrencyGateway();
