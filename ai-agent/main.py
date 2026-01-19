import json
from typing import Dict, Any
from langchain_ollama import ChatOllama
from langchain_core.prompts import ChatPromptTemplate
from pydantic import BaseModel, Field

# 1. Define the Data Models
class Node(BaseModel):
    id: str = Field(description="Unique identifier for the node. Use 'start' for the start node and 'end' for the end node.")
    type: str = Field(description="Type of the node. Options: 'Start', 'End', 'Move'.")
    params: Dict[str, Any] = Field(default_factory=dict, description="Parameters for the node (e.g., x, y, orientation).")
    transitions: Dict[str, str] = Field(default_factory=dict, description="Dictionary of transitions, e.g., {'next': 'target_node_id'}.")

class FlowResponse(BaseModel):
    """The full flow chart structure."""
    nodes: Dict[str, Node] = Field(description="Dictionary where keys are node IDs and values are Node objects.")

# 2. Setup the LLM
# Note: Ensure you have 'llama3' (or your preferred model) pulled in Ollama
# You can change the model name below if needed (e.g., "mistral", "llama2")
llm = ChatOllama(model="deepseek-r1:32b", temperature=0)

# 3. Create the Prompt
system_template = """You are an AI assistant that generates Robot Flow Charts for React Flow.
Your task is to convert natural language descriptions of robot movements into a structured JSON format.

The flow chart nodes are:
- Start: The entry point. id="start". Transition "next" points to the first action.
- End: The exit point. id="end". No transitions.
- Move: A movement command. Params: x (float), y (float), orientation (float). Transition "next" points to the next node.

Rules:
1. ALWAYS include a 'start' node and an 'end' node.
2. Generate unique IDs for intermediate nodes (e.g., timestamps or random strings).
3. Connect nodes logically using the 'transitions' field.
4. The output must be a valid JSON object where keys are node IDs.
5. Ensure that no movement command follows a sit command.
"""

prompt = ChatPromptTemplate.from_messages([
    ("system", system_template),
    ("human", "{input}"),
])

# 4. Chain with Structured Output
# We use with_structured_output to ensure the LLM matches our Schema
structured_llm = llm.with_structured_output(FlowResponse)
chain = prompt | structured_llm

def generate_flow(user_input: str):
    print(f"Generating flow for: '{user_input}'...")
    try:
        result = chain.invoke({"input": user_input})
        
        # Define a helper to convert Pydantic models to dicts for printing
        def pydantic_encoder(obj):
            if hasattr(obj, "dict"):
                return obj.dict()
            if hasattr(obj, "model_dump"):
                return obj.model_dump()
            return str(obj)

        # Print the 'nodes' part of the response nicely formatted
        print(json.dumps(result.nodes, default=pydantic_encoder, indent=2))
        return result.nodes
    except Exception as e:
        print(f"Error generating flow: {e}")
        # Fallback debug
        print("Raw result might not have matched schema. Check Ollama model capabilities.")

def main():
    print("=== Robot Flow Agent ===")
    print("Describe a movement (e.g., 'Move to x=5, y=5 then stop').")
    print("Type 'exit' to quit.\n")
    
    while True:
        try:
            user_input = input("User Input: ")
            if user_input.lower() in ["exit", "quit"]:
                break
            if not user_input.strip():
                continue
            
            generate_flow(user_input)
            print("-" * 50)
            
        except KeyboardInterrupt:
            print("\nExiting...")
            break

if __name__ == "__main__":
    main()
