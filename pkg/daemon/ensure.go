package daemon

import (
	"fmt"

	"github.com/0xSplits/indexingo/client"
	"github.com/0xSplits/indexingo/filters"
	"github.com/0xSplits/indexingo/pipelines"
	"github.com/0xSplits/indexingo/transformations"
	"github.com/xh3b4sd/tracer"
)

func (d *Daemon) Ensure() error {
	var cli client.Interface
	{
		cli = client.New(client.Config{
			Key: d.env.IndexingcoApiKey,
		})
	}

	//--------------------------------------------------------------------------//

	var fil *filters.Filters
	{
		fil = filters.New(filters.Config{
			Cli: cli,
		})
	}

	var tra *transformations.Transformations
	{
		tra = transformations.New(transformations.Config{
			Cli: cli,
		})
	}

	var pip *pipelines.Pipelines
	{
		pip = pipelines.New(pipelines.Config{
			Cli: cli,
		})
	}

	//--------------------------------------------------------------------------//

	{
		res, err := fil.AddValues("test-filter", []string{"0xb7f5bf799fb265657c628ef4a13f90f83a3a616a"})
		if err != nil {
			tracer.Panic(tracer.Mask(err))
		}

		d.log.Log(
			"level", "info",
			"message", "filter creation",
			"status", res.Message,
		)
	}

	//--------------------------------------------------------------------------//

	//
	//     curl -s --location --globoff 'https://app.indexing.co/dw/transformations/test?network=base&beat=37740907&filter=xh3b4sd-test-filter&filterKeys[0]=to&filterKeys[1]=from' --header "x-api-key: $INDEXINGCO_API_KEY" --form 'code="function traByBlock(blo) { const tra = templates.tokenTransfers(blo); return tra.map(x => ({ network: blo._network, chainId: utils.evmChainToId(blo._network), blockHash: blo.hash, blockNumber: blo.number, timestamp: utils.blockToTimestamp(blo), ...x, })); }"' | jq .
	//

	var cod string
	{
		cod = `
			function traByBlock(blo) {
				const tra = templates.tokenTransfers(blo);

				return tra.map(x => ({
					network: blo._network,
					chainId: utils.evmChainToId(blo._network),
					blockHash: blo.hash,
					blockNumber: blo.number,
					timestamp: utils.blockToTimestamp(blo),
					...x,
				}));
			}
		`
	}

	{
		res, err := tra.CreateTransformation("test-transformation", cod)
		if err != nil {
			tracer.Panic(tracer.Mask(err))
		}

		d.log.Log(
			"level", "info",
			"message", "transformation creation",
			"status", res.Message,
		)
	}

	//--------------------------------------------------------------------------//

	var cpr pipelines.CreatePipelineRequest
	{
		cpr = pipelines.CreatePipelineRequest{
			Name:           "test-pipeline",
			Transformation: "test-transformation",
			Filter:         "test-filter",
			FilterKeys:     []string{"from", "to"},
			Networks:       []string{"ethereum", "base"},
			Enabled:        true,
			Delivery: pipelines.CreatePipelineRequestDelivery{
				Adapter: "WEBSOCKET",
				Connection: pipelines.CreatePipelineRequestDeliveryConnection{
					Host: fmt.Sprintf("https://pulsar.%s.splits.org/indexing", d.env.Environment),
					Headers: map[string]string{
						"Authorization": "Bearer " + d.env.WebsocketSecret,
					},
				},
			},
		}
	}

	{
		res, err := pip.CreatePipeline(cpr)
		if err != nil {
			tracer.Panic(tracer.Mask(err))
		}

		d.log.Log(
			"level", "info",
			"message", "pipeline creation",
			"status", res.Message,
		)
	}

	//--------------------------------------------------------------------------//

	var bpr pipelines.BackfillPipelineRequest
	{
		bpr = pipelines.BackfillPipelineRequest{
			Network:   "base",
			Value:     "0xb7f5bf799fb265657c628ef4a13f90f83a3a616a",
			BeatStart: 37740907,
		}
	}

	{
		res, err := pip.BackfillPipeline("test-pipeline", bpr)
		if err != nil {
			tracer.Panic(tracer.Mask(err))
		}

		d.log.Log(
			"level", "info",
			"message", "backfill instruction",
			"status", res.Message,
		)
	}

	return nil
}
