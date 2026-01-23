import structlog
import uvicorn
from interfaces.rest.rest_server import init_rest_server
from agent.flow import init_flow_agent

def main():
    
    logger = structlog.get_logger()

    flow_agent = init_flow_agent(logger=logger)

    # register REST server
    app = init_rest_server(
        logger=logger,
        flow_agent=flow_agent
    )

    try:
        logger.info("Running API Server")
        uvicorn.run(app, host="0.0.0.0", port=3001)
    except KeyboardInterrupt as err:
        logger.error("[RUN] Uvicorn fun fastapi server failed: {}".format(err))
        exit()


if __name__ == "__main__":
    main()
