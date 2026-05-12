// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

import type { NextApiRequest, NextApiResponse } from 'next';
import InstrumentationMiddleware from '../../utils/telemetry/InstrumentationMiddleware';
import CurrencyGateway from '../../gateways/rpc/Currency.gateway';
import { Empty } from '../../protos/demo';

type TResponse = string[] | Empty;

const handler = async ({ method }: NextApiRequest, res: NextApiResponse<TResponse>) => {
  switch (method) {
    case 'GET': {
      try {
        const { currencyCodes = [] } = await CurrencyGateway.getSupportedCurrencies();

        return res.status(200).json(currencyCodes);
      } catch (error) {
        console.error('Failed to fetch currencies:', error);
        return res.status(503).json([]);
      }
    }

    default: {
      return res.status(405);
    }
  }
};

export default InstrumentationMiddleware(handler);
