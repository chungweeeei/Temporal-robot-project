import structlog
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from agent.flow import FlowAgent
from interfaces.rest.routers.flow_agent import init_agent_router

def init_rest_server(
    logger: structlog.stdlib.BoundLogger,
    flow_agent: FlowAgent,
) -> FastAPI:
    swagger_ui_desc = """
    This is the REST API documentation for the generate flow AI agent.
    """

    app = FastAPI(
        title="Generate Flow AI Agent", description=swagger_ui_desc, version="1.0.0"
    )
    app.add_middleware(
        CORSMiddleware,
        allow_origins=["*"],
        allow_credentials=True,
        allow_methods=["GET", "POST", "PUT", "DELETE", "OPTIONS"],
        allow_headers=["*"],
    )

    flow_router = init_agent_router(logger=logger, flow_agent=flow_agent)

    app.include_router(flow_router)

    return app
