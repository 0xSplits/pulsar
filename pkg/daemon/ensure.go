package daemon

import (
	"fmt"

	"github.com/0xSplits/indexingo/pipelines"
	"github.com/xh3b4sd/tracer"
)

func (d *Daemon) Ensure() error {
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
		res, err := d.tra.CreateTransformation(fmt.Sprintf("%s-transformation", d.env.Environment), cod)
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
			Name:           fmt.Sprintf("%s-pipeline", d.env.Environment),
			Transformation: fmt.Sprintf("%s-transformation", d.env.Environment),
			Filter:         fmt.Sprintf("%s-filter", d.env.Environment),
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
		res, err := d.pip.CreatePipeline(cpr)
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

	// var bpr pipelines.BackfillPipelineRequest
	// {
	// 	bpr = pipelines.BackfillPipelineRequest{
	// 		Network:   "base",
	// 		Value:     "0xb7f5bf799fb265657c628ef4a13f90f83a3a616a",
	// 		BeatStart: 37740907,
	// 	}
	// }

	// {
	// 	res, err := d.pip.BackfillPipeline(fmt.Sprintf("%s-pipeline", d.env.Environment), bpr)
	// 	if err != nil {
	// 		tracer.Panic(tracer.Mask(err))
	// 	}

	// 	d.log.Log(
	// 		"level", "info",
	// 		"message", "backfill instruction",
	// 		"status", res.Message,
	// 	)
	// }

	return nil
}
