import { useNavigate } from "react-router-dom";
import { Card, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { WorkflowStatusBadge } from "@/components/shared/WorkflowStatusBadge";
import type { WorkflowStatusDef } from "@/types/workflows";
import { Play, Pencil, Trash2, Pause } from "lucide-react";

interface WorkflowCardProps {
  workflowId: string;
  workflowName: string;
  status: WorkflowStatusDef;
  currentStep: string;
  onTrigger: (workflowId: string) => void;
  onDelete: (workflowId: string) => void;
  onPause: (workflowId: string) => void;
}

export function WorkflowCard({
  workflowId,
  workflowName,
  status,
  currentStep,
  onTrigger,
  onDelete,
  onPause,
}: WorkflowCardProps) {

  const navigate = useNavigate();

  const handleDelete = () => {
    if (window.confirm(`Are you sure you want to delete the workflow "${workflowName}"? This action cannot be undone.`)) {
      onDelete(workflowId);
    }
  };

  return (
    <Card className="hover:shadow-lg transition-shadow">
      <CardHeader className="pb-2">
          <div className="flex flex-row items-center justify-between gap-4">
            <CardTitle className="text-lg font-semibold truncate min-w-0" title={workflowName}>
              {workflowName || "Untitled Workflow"}
            </CardTitle>
            <div className="shrink-0">
              <WorkflowStatusBadge status={status} currentStep={currentStep} />
            </div>
          </div>
      </CardHeader>
      <CardFooter className="flex justify-between gap-2 pt-0">
        <div className="flex gap-1">
          <Button
            variant="outline"
            size="sm"
            onClick={() => navigate(`/editor/${workflowId}`, {
              state: {
                operation: "edit",
                workflowName: workflowName
              }
            })}
            disabled={false}
            title="Edit workflow"
          >
            <Pencil className="h-4 w-4" />
          </Button>
          <Button
            variant="outline"
            size="sm"
            onClick={handleDelete}
            disabled={false}
            className="text-destructive hover:text-destructive"
            title="Delete workflow"
          >
            <Trash2 className="h-4 w-4" />
          </Button>
        </div>
        <div>
          {(status === "Running" || status === "Paused") ? (
            <Button
              variant={status === "Paused" ? "default" : "secondary"}
              size="sm"
              onClick={() => onPause(workflowId)}
              title={status === "Paused" ? "Resume workflow" : "Pause workflow"}
            >
              <Pause className="h-4 w-4 mr-1" />
              {status === "Paused" ? "Resume" : "Pause"}
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
