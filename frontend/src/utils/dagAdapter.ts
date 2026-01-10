import type { Edge, Node } from '@xyflow/react';
import type { WorkflowPayload, WorkflowNode, FlowNodeData, WorkflowTransitions } from '../types/schema';

export const transformToDagPayload = (
  nodes: Node[],
  edges: Edge[]
): WorkflowPayload => {
  const payloadNodes: Record<string, WorkflowNode> = {};

  // 1. 初始化所有節點
  nodes.forEach((node) => {
    payloadNodes[node.id] = {
      id: node.id,
      type: (node.data as FlowNodeData).activityType,
      params: (node.data as FlowNodeData).params,
      transitions: {},
      retryPolicy: (node.data as FlowNodeData).retryPolicy,
    };
  });

  // 2. 根據 Edges 建立關聯
  edges.forEach((edge) => {
    const sourceNode = payloadNodes[edge.source];
    const targetId = edge.target;

    if (sourceNode) {
      // 根據 Handle ID 對應到正確的 transition slot
      // Action Node handles: 'success' (or null), 'failure'
      // Condition Node handles: 'true', 'false'
      
      const handleId = edge.sourceHandle;

      if (handleId === 'failure') {
        sourceNode.transitions.failure = targetId;
      } else if (handleId === 'true') {
        sourceNode.transitions.true = targetId;
      } else if (handleId === 'false') {
        sourceNode.transitions.false = targetId;
      } else {
        // default / success handle
        sourceNode.transitions.next = targetId;
      }
    }
  });

  // 3. 找出 Root Node (Start Node通常是id='start'，或者我們找沒有被指到的)
  // 這裡我們直接假設 id='start' 的是起點，或者 fallback 到計算入度
  let rootNodeId = 'start';
  
  // 如果沒有 id 為 'start' 的 node，則透過 In-Degree 計算
  if (!payloadNodes['start']) {
    const inDegree: Record<string, number> = {};
    Object.keys(payloadNodes).forEach(id => inDegree[id] = 0);
    edges.forEach(edge => {
       inDegree[edge.target] = (inDegree[edge.target] || 0) + 1;
    });
    const roots = Object.keys(inDegree).filter(id => inDegree[id] === 0);
    if (roots.length > 0) rootNodeId = roots[0];
  }

  return {
    workflowId: crypto.randomUUID(), // Generate a random ID for this execution/template
    rootNodeId,
    nodes: payloadNodes,
  };
};
