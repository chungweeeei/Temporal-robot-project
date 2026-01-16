import { Badge } from "@/components/ui/badge";
import type { WorkflowStatus } from "@/types/schema";

interface WorkflowStatusBadgeProps {
  status: WorkflowStatus;
  currentStep?: string;
}

const statusConfig: Record<WorkflowStatus, { label: string; variant: "default" | "secondary" | "destructive" | "outline"; className: string }> = {
  idle: {
    label: "Idle",
    variant: "secondary",
    className: "bg-gray-100 text-gray-700 hover:bg-gray-100",
  },
  running: {
    label: "Running",
    variant: "default",
    className: "bg-blue-500 text-white hover:bg-blue-500 animate-pulse",
  },
  paused: {
    label: "Paused",
    variant: "outline",
    className: "border-amber-500 text-amber-600 hover:bg-amber-50",
  },
  completed: {
    label: "Completed",
    variant: "default",
    className: "bg-green-500 text-white hover:bg-green-500",
  },
  failed: {
    label: "Failed",
    variant: "destructive",
    className: "bg-red-500 text-white hover:bg-red-500",
  },
};

export function WorkflowStatusBadge({ status, currentStep }: WorkflowStatusBadgeProps) {
  const config = statusConfig[status];

  return (
    <div className="flex items-center gap-2">
      <Badge variant={config.variant} className={config.className}>
        {config.label}
      </Badge>
      {status === "running" && currentStep && (
        <span className="text-xs text-muted-foreground">
          Step: <span className="font-medium">{currentStep}</span>
        </span>
      )}
    </div>
  );
}
