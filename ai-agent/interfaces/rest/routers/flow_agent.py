import structlog
from fastapi import APIRouter, HTTPException, status

from agent.flow import FlowAgent

def init_agent_router(
    logger: structlog.stdlib.BoundLogger,
    flow_agent: FlowAgent
) -> APIRouter:
    
    router = APIRouter(prefix="", tags=["agent"])

    @router.post("/api/v1/agent/flows/generate")
    def generate_flow(message: str):
        try:
            flow = flow_agent.generate_flow(message=message)
        except Exception as e:
            logger.error(f"Flow generation failed: {e}")
            raise HTTPException(
                status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
                detail="Flow generation failed"
            )
        return flow

    return router

