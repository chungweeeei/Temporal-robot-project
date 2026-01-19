import { useState, useCallback, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useLocation } from 'react-router-dom';
import { ReactFlow, Background, Controls, type Node, type Edge, addEdge, type OnNodesChange, type OnEdgesChange, type OnNodesDelete, applyNodeChanges, applyEdgeChanges, type NodeMouseHandler, type Connection } from '@xyflow/react';
import '@xyflow/react/dist/style.css';

import type { BaseParams, MoveParams, SleepParams } from '@/types/schema';
import ConditionNode from '@/components/nodes/ConditionNode';
import ActionNode from '@/components/nodes/ActionNode';
import StartNode from '@/components/nodes/StartNode';
import EndNode from '@/components/nodes/EndNode';
import NodeEditorModal from '@/components/NodeEditorModal';
import { WorkflowStatusBadge } from '@/components/shared/WorkflowStatusBadge';
import { transformBackToReactFlow, transformToDagPayload } from '@/utils/dagAdapter';
import { useSaveWorkflow } from '@/hooks/useSaveWorkflow';
import { useTriggerWorkflow } from '@/hooks/useTriggerWorkflow';
import { useFetchWorkflowById } from '@/hooks/useFetchWorkflowById';
import { useWorkflowMonitor } from '@/hooks/useWorkflowMonitor';
import { usePauseWorkflow, useResumeWorkflow } from '@/hooks/useControlWorkflow';
import { Button } from '@/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { ArrowLeft, Save, Play, Pause, Plus, ChevronDown } from 'lucide-react';

// --- Register Custom Nodes ---
const nodeTypes = {
  condition: ConditionNode,
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

export default function Editor() {
  const { workflowId } = useParams<{ workflowId: string }>();

  // state passed by useLocation in react-router
  const location = useLocation();
  const state = location.state as { operation?: string; workflowName?: string };
  const isCreateMode = state?.operation === 'create';
  
  // React Router navigation
  const navigate = useNavigate();
  
  // --- React Flow State ---
  const [nodes, setNodes] = useState<Node[]>([]);
  const [edges, setEdges] = useState<Edge[]>([]);

  // Fetch workflow detail from backend
  const { data: workflowDetail, isLoading } = useFetchWorkflowById(workflowId || '', { enabled: !isCreateMode});
  const workflowName = workflowDetail?.workflow_name || state.workflowName;

  // Effect: When workflow detail is loaded, update the canvas
  useEffect(() => {
    if (!workflowDetail || !workflowDetail.nodes){
      // set default start and end nodes for new workflow
      setNodes([
        {
          id: 'start',
          type: 'start',
          position: { x: 0, y: 0 },
          data: { label: 'Start' },
        },
        {
          id: 'end',
          type: 'end',
          position: { x: 1000, y: 0 },
          data: { label: 'End' },
        }
      ]);
      setEdges([]);
      return;
    }
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
    const payload = transformToDagPayload(workflowId, workflowName, nodes, edges);
    saveMutation.mutate(payload, {
      onSuccess: () => {
        alert('Workflow saved successfully!');
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
    const payload = transformToDagPayload(workflowId, workflowName, nodes, edges);
    triggerMutation.mutate(payload, {
      onSuccess: () => {
        alert('Workflow triggered successfully!');
        setIsMonitoring(true);
      },
      onError: (error) => {
        alert(`Failed to trigger: ${error.message}`);
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
      type: isCondition ? 'condition' : 'action',
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

  if (isLoading) {
    return (
      <div className="w-screen h-screen flex items-center justify-center">
        <div className="text-muted-foreground">Loading workflow...</div>
      </div>
    );
  }

  const isRunning = workflowStatus === 'running';
  const isPaused = workflowStatus === 'paused';

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
              <WorkflowStatusBadge status={workflowStatus} currentStep={currentStep} />
            </div>
          </div>

          {/* Center: Add Node Buttons */}
          <div className="flex items-center gap-2">
            {activityOptions.map((option) => (
              <Button
                key={option.type}
                variant="outline"
                size="sm"
                onClick={() => handleAddNode(option.type)}
                disabled={isRunning}
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
              disabled={saveMutation.isPending || isRunning}
            >
              <Save className="h-4 w-4 mr-2" />
              {saveMutation.isPending ? 'Saving...' : 'Save'}
            </Button>

            {isRunning || isPaused ? (
              <Button
                variant={isPaused ? 'default' : 'secondary'}
                size="sm"
                onClick={isPaused ? handleResumeWorkflow : handlePauseWorkflow}
              >
                <Pause className="h-4 w-4 mr-2" />
                {isPaused ? 'Resume' : 'Pause'}
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

        {/* Tablet/Mobile Layout */}
        <div className="lg:hidden space-y-3">
          {/* Row 1: Back + Title + Actions */}
          <div className="flex items-center justify-between">
            {/* Left: Back + Title */}
            <div className="flex items-center gap-2 sm:gap-4 min-w-0 flex-1">
              <Button variant="ghost" size="sm" onClick={() => navigate('/')} className="shrink-0">
                <ArrowLeft className="h-4 w-4 sm:mr-2" />
                <span className="hidden sm:inline">Back</span>
              </Button>
              <div className="h-6 w-px bg-border shrink-0" />
              <div className="min-w-0">
                <h1 className="font-semibold truncate">{workflowName}</h1>
                <WorkflowStatusBadge status={workflowStatus} currentStep={currentStep} />
              </div>
            </div>

            {/* Right: Actions */}
            <div className="flex items-center gap-2 shrink-0">
              <Button
                variant="outline"
                size="sm"
                onClick={handleSaveWorkflow}
                disabled={saveMutation.isPending || isRunning}
              >
                <Save className="h-4 w-4 sm:mr-2" />
                <span className="hidden sm:inline">{saveMutation.isPending ? 'Saving...' : 'Save'}</span>
              </Button>

              {isRunning || isPaused ? (
                <Button
                  variant={isPaused ? 'default' : 'secondary'}
                  size="sm"
                  onClick={isPaused ? handleResumeWorkflow : handlePauseWorkflow}
                >
                  <Pause className="h-4 w-4 sm:mr-2" />
                  <span className="hidden sm:inline">{isPaused ? 'Resume' : 'Pause'}</span>
                </Button>
              ) : (
                <Button
                  variant="default"
                  size="sm"
                  onClick={handleTriggerWorkflow}
                  disabled={triggerMutation.isPending}
                >
                  <Play className="h-4 w-4 sm:mr-2" />
                  <span className="hidden sm:inline">{triggerMutation.isPending ? 'Starting...' : 'Run'}</span>
                </Button>
              )}
            </div>
          </div>

          {/* Row 2: Add Node - Dropdown on mobile, buttons on tablet */}
          <div className="flex items-center gap-2">
            {/* Mobile: Dropdown */}
            <div className="sm:hidden">
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button variant="outline" size="sm" disabled={isRunning}>
                    <Plus className="h-4 w-4 mr-2" />
                    Add Node
                    <ChevronDown className="h-4 w-4 ml-2" />
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="start">
                  {activityOptions.map((option) => (
                    <DropdownMenuItem
                      key={option.type}
                      onClick={() => handleAddNode(option.type)}
                    >
                      <Plus className="h-4 w-4 mr-2" />
                      {option.label}
                    </DropdownMenuItem>
                  ))}
                </DropdownMenuContent>
              </DropdownMenu>
            </div>

            {/* Tablet: Scrollable buttons */}
            <div className="hidden sm:flex items-center gap-2 overflow-x-auto pb-1">
              {activityOptions.map((option) => (
                <Button
                  key={option.type}
                  variant="outline"
                  size="sm"
                  onClick={() => handleAddNode(option.type)}
                  disabled={isRunning}
                  className="shrink-0"
                >
                  <Plus className="h-3 w-3 mr-1" />
                  {option.label}
                </Button>
              ))}
            </div>
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
