import { useNavigate } from "react-router-dom";
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { WorkflowStatusBadge } from "@/components/shared/WorkflowStatusBadge";
import type { WorkflowStatus } from "@/types/schema";
import { Play, Pencil, Trash2, Pause } from "lucide-react";

interface WorkflowCardProps {
  workflowId: string;
  workflowName: string;
  status: WorkflowStatus;
  currentStep?: string;
  onTrigger: (id: string) => void;
  onPause: (id: string) => void;
  onDelete: (id: string) => void;
}

export function WorkflowCard({
  workflowId,
  workflowName,
  status,
  currentStep,
  onTrigger,
  onPause,
  onDelete,
}: WorkflowCardProps) {
  const navigate = useNavigate();

  const handleEdit = () => {
    navigate(`/editor/${workflowId}`);
  };

  const handleDelete = () => {
    if (window.confirm(`確定要刪除 "${workflowName}" 嗎？此操作無法復原。`)) {
      onDelete(workflowId);
    }
  };

  const isRunning = status === "running";
  const isPaused = status === "paused";

  return (
    <Card className="hover:shadow-lg transition-shadow">
      <CardHeader className="pb-2">
        <CardTitle className="text-lg font-semibold truncate" title={workflowName}>
          {workflowName || "Untitled Workflow"}
        </CardTitle>
      </CardHeader>
      <CardContent className="pb-3">
        <WorkflowStatusBadge status={status} currentStep={currentStep} />
      </CardContent>
      <CardFooter className="flex justify-between gap-2 pt-0">
        <div className="flex gap-1">
          <Button
            variant="outline"
            size="sm"
            onClick={handleEdit}
            disabled={isRunning}
            title="Edit workflow"
          >
            <Pencil className="h-4 w-4" />
          </Button>
          <Button
            variant="outline"
            size="sm"
            onClick={handleDelete}
            disabled={isRunning}
            className="text-destructive hover:text-destructive"
            title="Delete workflow"
          >
            <Trash2 className="h-4 w-4" />
          </Button>
        </div>
        <div>
          {isRunning || isPaused ? (
            <Button
              variant={isPaused ? "default" : "secondary"}
              size="sm"
              onClick={() => onPause(workflowId)}
              title={isPaused ? "Resume workflow" : "Pause workflow"}
            >
              <Pause className="h-4 w-4 mr-1" />
              {isPaused ? "Resume" : "Pause"}
            </Button>
          ) : (
            <Button
              variant="default"
              size="sm"
              onClick={() => onTrigger(workflowId)}
              title="Run workflow"
            >
              <Play className="h-4 w-4 mr-1" />
              Run
            </Button>
          )}
        </div>
      </CardFooter>
    </Card>
  );
}
