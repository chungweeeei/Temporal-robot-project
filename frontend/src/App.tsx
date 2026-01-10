import { useState, useCallback } from 'react';
import { ReactFlow, Background, Controls, type Node, type Edge, addEdge, type OnNodesChange, type OnEdgesChange, type OnNodesDelete, applyNodeChanges, applyEdgeChanges, type NodeMouseHandler, type Connection } from '@xyflow/react';
import '@xyflow/react/dist/style.css';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

// 引入我們剛剛定義的型別與工具
import { transformToDagPayload } from './utils/dagAdapter';
import type { BaseParams, RetryPolicy, ConditionParams, MoveParams, SleepParams } from './types/schema';
import ConditionNode from './components/nodes/ConditionNode';
import ActionNode from './components/nodes/ActionNode';
import StartNode from './components/nodes/StartNode';
import EndNode from './components/nodes/EndNode';
import WorkflowToolbar, { type WorkflowStatus } from './components/WorkflowToolbar';
import NodeEditorModal from './components/NodeEditorModal';

// --- 註冊 Custom Nodes ---
const nodeTypes = {
  condition: ConditionNode,
  action: ActionNode,
  start: StartNode,
  end: EndNode,
};

// --- 建立 React Query Client ---
const queryClient = new QueryClient();

// --- 模擬初始節點 ---
const initialNodes: Node[] = [
  { id: 'start', position: { x: 50, y: 300 }, data: { label: 'Start', activityType: 'Start', params: {} }, type: 'start', deletable: false },
  { id: 'end', position: { x: 600, y: 300 }, data: { label: 'End', activityType: 'End', params: {} }, type: 'end', deletable: false },
];

function Scheduler() {
  const [nodes, setNodes] = useState<Node[]>(initialNodes);
  const [edges, setEdges] = useState<Edge[]>([]);
  
  // --- Modal 狀態 ---
  const [editingNode, setEditingNode] = useState<Node | null>(null);
  const [modalOpen, setModalOpen] = useState(false);

  // --- Workflow Execution State ---
  const [workflowStatus, setWorkflowStatus] = useState<WorkflowStatus>('idle');

  // React Flow 回呼函式
  const onNodesChange: OnNodesChange = useCallback((changes) => setNodes((nds) => applyNodeChanges(changes, nds)), []);
  const onEdgesChange: OnEdgesChange = useCallback((changes) => setEdges((eds) => applyEdgeChanges(changes, eds)), []);
  const onConnect = useCallback((params: Connection) => setEdges((eds) => addEdge(params, eds)), []);

  const onNodesDelete: OnNodesDelete = useCallback((deleted) => {
    setEdges((eds) => {
      const deletedIds = new Set(deleted.map((node) => node.id));
      return eds.filter((edge) => !deletedIds.has(edge.source) && !deletedIds.has(edge.target));
    });
  }, []);

  // --- 核心：雙擊節點開啟編輯 ---
  const onNodeDoubleClick: NodeMouseHandler = useCallback((_event, node) => {
    if (node.id === 'start' || node.id === 'end') return;
    setEditingNode(node);
    setModalOpen(true);
  }, []);

  // --- 儲存 Modal 的變更 ---
  const handleSaveNodeParams = (newParams: BaseParams | MoveParams | SleepParams | ConditionParams, newRetryPolicy?: RetryPolicy) => {
    if (!editingNode) return;
    
    setNodes((nds) => nds.map((node) => {
      if (node.id === editingNode.id) {
        // 更新節點的 data (UI 會自動更新)
        return { 
          ...node, 
          data: { 
            ...node.data, 
            params: newParams,
            retryPolicy: newRetryPolicy || node.data.retryPolicy
          } 
        };
      }
      return node;
    }));
    setModalOpen(false);
  };

  // --- API Mutation (儲存 Workflow) ---
  /* const mutation = useMutation({
    mutationFn: (data: WorkflowPayload) => axios.post('http://localhost:3000/api/workflow', data),
    onSuccess: () => alert('Workflow saved successfully!'),
    onError: (err) => alert(`Error: ${err.message}`)
  }); */

  const handleSaveWorkflow = () => {
    const payload = transformToDagPayload(nodes, edges);
    console.log("Sending Payload:", JSON.stringify(payload, null, 2));
    // mutation.mutate(payload); // 暫時註解，待後端完成
    alert('Payload generated! Check console for details.');
  };

  const handleTriggerWorkflow = () => {
    // Placeholder for trigger logic
    console.log("Trigger workflow clicked");
    setWorkflowStatus('running');
    
    // Simulate a workflow completion after 5 seconds
    setTimeout(() => {
      setWorkflowStatus((prev) => prev === 'running' ? 'completed' : prev);
    }, 5000);
  };

  const handleStopWorkflow = () => {
    console.log("Stop workflow clicked");
    setWorkflowStatus('paused');
  };

  const handleResumeWorkflow = () => {
    console.log("Resume workflow clicked");
    if (workflowStatus === 'paused') {
      setWorkflowStatus('running');
    }
  };

  // --- 新增節點功能 ---
  const handleAddNode = (type: string) => {
    const id = Date.now().toString();
    const isCondition = type === 'Condition';
    
    const newNode: Node = {
      id,
      position: { x: Math.random() * 400, y: Math.random() * 400 },
      data: { 
        label: type,
        activityType: type, 
        params: type === 'Move' ? { x: 0, y: 0 } 
              : type === 'Sleep' ? { duration: 1000 } 
              : type === 'Condition' ? { expression: 'true' }
              : {},
        retryPolicy: {
          maxAttempts: 1,
          initialInterval: 1000,
        }
      },
      type: isCondition ? 'condition' : 'action', // use 'action' instead of 'default'
    };
    setNodes((nds) => [...nds, newNode]);
  };

  // --- 刪除節點功能 ---
  const handleDeleteNode = () => {
    if (!editingNode) return;
    setNodes((nds) => nds.filter((n) => n.id !== editingNode.id));
    setEdges((eds) => eds.filter((e) => e.source !== editingNode.id && e.target !== editingNode.id));
    setModalOpen(false);
    setEditingNode(null);
  };

  return (
    <div className="w-screen h-screen flex flex-col">
      <WorkflowToolbar 
        onAddNode={handleAddNode} 
        onSave={handleSaveWorkflow} 
        onTrigger={handleTriggerWorkflow}
        onStop={handleStopWorkflow}
        onResume={handleResumeWorkflow}
        status={workflowStatus}
      />

      <div className="flex-1">
        <ReactFlow
          nodes={nodes}
          edges={edges}
          onNodesChange={onNodesChange}
          onEdgesChange={onEdgesChange}
          onConnect={onConnect}
          onNodesDelete={onNodesDelete}
          onNodeDoubleClick={onNodeDoubleClick}
          nodeTypes={nodeTypes}
          fitView
        >
          <Background />
          <Controls />
        </ReactFlow>
      </div>

      <NodeEditorModal
        isOpen={modalOpen}
        node={editingNode}
        onClose={() => setModalOpen(false)}
        onSave={handleSaveNodeParams}
        onDelete={handleDeleteNode}
      />
    </div>
  );
}

// --- 包覆 Provider ---
export default function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <Scheduler />
    </QueryClientProvider>
  );
}
