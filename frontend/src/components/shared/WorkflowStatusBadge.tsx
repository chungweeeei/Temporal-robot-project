import { Badge } from "@/components/ui/badge";
import type { WorkflowStatusDef } from "@/types/workflows";

interface WorkflowStatusBadgeProps {
  status: WorkflowStatusDef;
  currentStep?: string;
}

const statusConfig: Record<WorkflowStatusDef, { label: string; variant: "default" | "secondary" | "destructive" | "outline"; className: string }> = {
  Idle: {
    label: "Idle",
    variant: "secondary",
    className: "text-sm bg-gray-100 text-gray-700 hover:bg-gray-100",
  },
  Running: {
    label: "Running",
    variant: "default",
    className: "text-sm bg-blue-500 text-white hover:bg-blue-500 animate-pulse",
  },
  Paused: {
    label: "Paused",
    variant: "outline",
    className: "text-sm border-amber-500 text-amber-600 hover:bg-amber-50",
  },
  Completed: {
    label: "Completed",
    variant: "default",
    className: "text-sm bg-green-500 text-white hover:bg-green-500",
  },
  Failed: {
    label: "Failed",
    variant: "destructive",
    className: "text-sm bg-red-500 text-white hover:bg-red-500",
  },
  Cancelled: {
    label: "Cancelled",
    variant: "outline",
    className: "text-sm border-gray-500 text-gray-600 hover:bg-gray-50",
  },
  Terminated: {
    label: "Terminated",
    variant: "destructive",
    className: "text-sm bg-orange-500 text-white hover:bg-orange-500",
  },
};

export function WorkflowStatusBadge({ status, currentStep }: WorkflowStatusBadgeProps) {
  const config = statusConfig[status];
  return (
    <div className="flex flex-col items-center gap-2">
      <Badge variant={config.variant} className={config.className}>
        {config.label}
      </Badge>
      {status === "Running" && currentStep && (
        <span className="text-sm text-muted-foreground">
          Step: <span className="font-medium">{currentStep}</span>
        </span>
      )}
    </div>
  );
}
