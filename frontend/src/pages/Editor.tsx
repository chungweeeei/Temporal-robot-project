import { useState, useCallback, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useLocation } from 'react-router-dom';
import { ReactFlow, Background, Controls, type Node, type Edge, addEdge, type OnNodesChange, type OnEdgesChange, type OnNodesDelete, applyNodeChanges, applyEdgeChanges, type NodeMouseHandler, type Connection } from '@xyflow/react';
import '@xyflow/react/dist/style.css';

import type { BaseParams, MoveParams, SleepParams } from '@/types/workflows';
import ActionNode from '@/components/nodes/ActionNode';
import StartNode from '@/components/nodes/StartNode';
import EndNode from '@/components/nodes/EndNode';
import NodeEditorModal from '@/components/NodeEditorModal';
import { WorkflowStatusBadge } from '@/components/shared/WorkflowStatusBadge';
import { transformBackToReactFlow, transformToDagPayload } from '@/utils/dagAdapter';
import { useSaveWorkflow } from '@/hooks/useSaveWorkflow';
import { useFetchWorkflowById } from '@/hooks/useFetchWorkflows';
import { useWorkflowMonitor } from '@/hooks/useWorkflowMonitor';
import { useTriggerWorkflow, usePauseWorkflow, useResumeWorkflow } from '@/hooks/useControlWorkflow';
import { Button } from '@/components/ui/button';
import { ArrowLeft, Save, Plus, Pause, Play } from 'lucide-react';

// --- Register Custom Nodes ---
const nodeTypes = {
  action: ActionNode,
  start: StartNode,
  end: EndNode,
};

// Activity type options
const activityOptions = [
  { type: 'Move', label: 'Move' },
  { type: 'Sleep', label: 'Sleep' },
  { type: 'Standup', label: 'Standup' },
  { type: 'Sitdown', label: 'Sitdown' },
  { type: 'TTS', label: 'TTS' },
  { type: 'Head', label: 'Head' }
];

const defaultNodes = [{
  id: 'start',
  type: 'start',
  position: { x: 0, y: 0 },
  data: { label: 'Start' },
},{
  id: 'end',
  type: 'end',
  position: { x: 1000, y: 0 },
  data: { label: 'End' },
}];

export default function Editor() {
  const { workflowId } = useParams<{ workflowId: string }>();

  // state passed by useLocation in react-router
  const location = useLocation();
  const state = location.state as { operation: string; workflowName?: string };
  const isCreateMode = state.operation === 'create';
  
  // Fetch workflow detail from backend
  const { data: workflowDetail, isLoading } = useFetchWorkflowById(workflowId || '', { enabled: !isCreateMode});
  const workflowName = workflowDetail?.workflow_name || state.workflowName;

  // React Router navigation
  const navigate = useNavigate();
  
  // --- React Flow State ---
  const [nodes, setNodes] = useState<Node[]>(defaultNodes);
  const [edges, setEdges] = useState<Edge[]>([]);

  // Effect: When workflow detail is loaded, update the canvas
  useEffect(() => {
    if (!workflowDetail || !workflowDetail.nodes) return;

    const { nodes: newNodes, edges: newEdges } = transformBackToReactFlow(workflowDetail.nodes);
    setNodes(newNodes);
    setEdges(newEdges);
  }, [workflowDetail]);

  // --- React Flow Node Modal state ---
  const [editingNode, setEditingNode] = useState<Node | null>(null);
  const [modalOpen, setModalOpen] = useState(false);

  // React Flow callback function
  const onNodesChange: OnNodesChange = useCallback((changes) => setNodes((nds) => applyNodeChanges(changes, nds)), []);
  const onEdgesChange: OnEdgesChange = useCallback((changes) => setEdges((eds) => applyEdgeChanges(changes, eds)), []);
  const onConnect = useCallback((params: Connection) => setEdges((eds) => addEdge(params, eds)), []);
  const onNodesDelete: OnNodesDelete = useCallback((deleted) => {
    setEdges((eds) => {
      const deletedIds = new Set(deleted.map((node) => node.id));
      return eds.filter((edge) => !deletedIds.has(edge.source) && !deletedIds.has(edge.target));
    });
  }, []);

  // --- Core: Double-click node to open editor ---
  const onNodeDoubleClick: NodeMouseHandler = useCallback((_event, node) => {
    if (node.id === 'start' || node.id === 'end') return;
    setEditingNode(node);
    setModalOpen(true);
  }, []);

  // --- Save changes from Modal ---
  const handleSaveNodeParams = (newParams: BaseParams | MoveParams | SleepParams) => {
    if (!editingNode) return;
    
    setNodes((nds) => nds.map((node) => {
      if (node.id === editingNode.id) {
        return { 
          ...node, 
          data: { 
            ...node.data, 
            params: newParams,
          } 
        };
      }
      return node;
    }));
    setModalOpen(false);
  };

  // --- API Mutation (儲存 Workflow) ---
  const saveMutation = useSaveWorkflow();
  const handleSaveWorkflow = () => {
    if (!workflowId) return;
    const payload = transformToDagPayload(nodes, edges);
    saveMutation.mutate({
      workflow_id: workflowId,
      workflow_name: workflowName,
      nodes: payload,
    }, {
      onSuccess: () => {
        alert("Workflow saved successfully.");
      },
      onError: (error) => {
        alert(`Failed to save: ${error.message}`);
      }
    });
  };

  // --- Workflow 執行狀態監控 ---
  const { setIsMonitoring, workflowStatus, currentStep } = useWorkflowMonitor(workflowId!);
  
  const triggerMutation = useTriggerWorkflow();
  const handleTriggerWorkflow = () => {
    if (!workflowId) return;
    triggerMutation.mutate(workflowId, {
      onSuccess: () => {
        console.log("Workflow triggered successfully.");
        setIsMonitoring(true);
      },
      onError: (error) => {
        console.log(`Failed to trigger: ${error.message}`);
        setIsMonitoring(false);
      }
    });
  };

  const pauseWorkflow = usePauseWorkflow();
  const handlePauseWorkflow = () => {
    if (!workflowId) return;
    pauseWorkflow.mutate(workflowId, {
      onSuccess: () => {
        alert('Workflow paused successfully!');
      },
      onError: (error) => {
        alert(`Failed to pause: ${error.message}`);
      }
    });
  };

  const resumeWorkflow = useResumeWorkflow();
  const handleResumeWorkflow = () => {
    if (!workflowId) return;
    resumeWorkflow.mutate(workflowId, {
      onSuccess: () => {
        alert('Workflow resumed successfully!');
      },
      onError: (error) => {
        alert(`Failed to resume: ${error.message}`);
      }
    });
  };

  // --- 新增節點功能 ---
  const handleAddNode = useCallback((type: string) => {
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
              : {},
        retryPolicy: {
          maxAttempts: 1,
          initialInterval: 1000,
        }
      },
      type: isCondition ? 'condition' : 'action',
    };
    setNodes((nds) => [...nds, newNode]);
  }, []);

  // --- 刪除節點功能 ---
  const handleDeleteNode = () => {
    if (!editingNode) return;
    setNodes((nds) => nds.filter((n) => n.id !== editingNode.id));
    setEdges((eds) => eds.filter((e) => e.source !== editingNode.id && e.target !== editingNode.id));
    setModalOpen(false);
    setEditingNode(null);
  };

  if (isLoading) {
    return (
      <div className="w-screen h-screen flex items-center justify-center">
        <div className="text-muted-foreground">Loading workflow...</div>
      </div>
    );
  }

  const isRunning = workflowStatus === 'Running';
  const isPaused = workflowStatus === 'Paused';

  return (
    <div className="w-screen h-screen flex flex-col">
      {/* Toolbar */}
      <header className="border-b bg-card px-4 py-3">
        {/* Desktop Layout */}
        <div className="hidden lg:flex items-center justify-between">
          {/* Left: Back + Title */}
          <div className="flex items-center gap-4">
            <Button variant="ghost" size="sm" onClick={() => navigate('/')}>
              <ArrowLeft className="h-4 w-4 mr-2" />
              Back
            </Button>
            <div className="h-6 w-px bg-border" />
            <div>
              <h1 className="font-semibold">{workflowName}</h1>
              <WorkflowStatusBadge status={workflowStatus} currentStep={currentStep || ""} />
            </div>
          </div>

          {/* 
            Center: Add Node Buttons 
            TODO: Action Node should fetch from backend
          */}
          <div className="flex items-center gap-2">
            {activityOptions.map((option) => (
              <Button
                key={option.type}
                variant="outline"
                size="sm"
                onClick={() => handleAddNode(option.type)}
                disabled={false}
              >
                <Plus className="h-3 w-3 mr-1" />
                {option.label}
              </Button>
            ))}
          </div>

          {/* Right: Actions */}
          <div className="flex items-center gap-2">
            <Button
              variant="outline"
              size="sm"
              onClick={handleSaveWorkflow}
            >
              <Save className="h-4 w-4 mr-2" />
              Save
            </Button>

            {isRunning || isPaused ? (
              <Button
                variant={isPaused ? 'default' : 'secondary'}
                size="sm"
                onClick={isPaused ? handleResumeWorkflow : handlePauseWorkflow}
              >
                {isPaused ? (
                  <>
                    <Play className="h-4 w-4 mr-2" />
                    Resume
                  </>
                  ) : (
                  <>
                    <Pause className="h-4 w-4 mr-2" />
                    Pause
                  </>
                )}
              </Button>
            ) : (
              <Button
                variant="default"
                size="sm"
                onClick={handleTriggerWorkflow}
                disabled={triggerMutation.isPending}
              >
                <Play className="h-4 w-4 mr-2" />
                {triggerMutation.isPending ? 'Starting...' : 'Run'}
              </Button>
            )}
          </div>
        </div>
      </header>

      {/* React Flow Canvas */}
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

      {/* Node Editor Modal */}
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
