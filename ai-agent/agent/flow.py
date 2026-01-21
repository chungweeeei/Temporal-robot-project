import json
from dotenv import load_dotenv
load_dotenv()

import structlog
from typing import Dict, Any
from langchain_ollama import ChatOllama
from langchain_core.prompts import ChatPromptTemplate

from pydantic import BaseModel, Field

class Node(BaseModel):
    id: str = Field(description="Unique identifier for the node. Use 'start' for the start node and 'end' for the end node.")
    type: str = Field(description="Type of the node. Options: 'Start', 'End', 'Move', 'Standup', 'Sitdown', 'TTS', 'Head'.")
    params: Dict[str, Any] = Field(default_factory=dict, description="Parameters for the node (e.g., x, y, orientation).")
    transitions: Dict[str, str] = Field(default_factory=dict, description="Dictionary of transitions, e.g., {'next': 'target_node_id'}.")

class FlowResponse(BaseModel):
    nodes: Dict[str, Node]

class FlowAgent:

    def __init__(self, logger: structlog.stdlib.BoundLogger):
        self.logger = logger

    def setup_llm(self, model_name: str):
        base_llm = ChatOllama(model=model_name)
        self._llm = base_llm.with_structured_output(FlowResponse)

    def setup_prompt(self):

        self._system_template = """
You are an AI assistant that generates Robot Flow Charts for React Flow.
Your task is to convert natural language descriptions of robot movements into a structured JSON format.

The flow chart nodes are included below:
- Start: The entry point. id="start". Transition "next" points to the first action.
- End: The exit point. id="end". No transitions.
- Move: A movement command. Params: x (float), y (float), orientation (float). Transition "next" points to the next node.
- Standup: Command to stand up. No params. Transition "next" points to the next node.
- Sitdown: Command to sit down. No params. Transition "next" points to the next node.
- TTS: Command to speak a message. Params: text (string). Transition "next" points to the next node.
- Head: Command to move head. Params: angle(float). Transition "next" points to the next node.

Rules:
1. ALWAYS include a 'start' node and an 'end' node.
2. Generate unique IDs for intermediate nodes (e.g., timestamps or random strings).
3. Connect nodes logically using the 'transitions' field.
4. The output must be a valid JSON object where keys are node IDs.
5. Ensure that no movement command follows a sit command.
6. Currently system do not support start connect to two different nodes.

Examples Output:
Examples Output:
{{
  "end": {{
    "id": "end",
    "type": "End",
    "params": {{}},
    "transitions": {{}}
  }},
  "start": {{
    "id": "start",
    "type": "Start",
    "params": {{}},
    "transitions": {{
      "next": "1768894485544"
    }}
  }},
  "1768894485544": {{
    "id": "1768894485544",
    "type": "Move",
    "params": {{
      "x": 3,
      "y": 3,
      "orientation": 0
    }},
    "transitions": {{
      "next": "1768894486397"
    }}
  }},
  "1768894486397": {{
    "id": "1768894486397",
    "type": "Sleep",
    "params": {{
      "duration": 5000
    }},
    "transitions": {{
      "next": "end"
    }}
  }}
}}
"""  

    def generate_flow(self, message: str) -> FlowResponse:
        prompt = ChatPromptTemplate.from_messages([
            ("system", self._system_template),
            ("user", "Generate a flow chart for the following request: {message}")
        ])

        chain = prompt | self._llm

        try:
            result = chain.invoke({"message": message})
            return result
        except Exception as e:
            self.logger.error(f"LLM invocation failed: {e}")
            raise


def init_flow_agent(logger: structlog.stdlib.BoundLogger) -> FlowAgent:
    flow_agent = FlowAgent(logger=logger)
    flow_agent.setup_llm(model_name="deepseek-r1:32b")
    flow_agent.setup_prompt()
    return flow_agent