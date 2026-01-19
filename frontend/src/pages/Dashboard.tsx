import { WorkflowList } from "@/components/dashboard/WorkflowList";
import { CreateWorkflowModal } from "@/components/dashboard/CreateWorkflowModal";
import { ScheduleSection } from "@/components/dashboard/ScheduleSection";
import { Bot } from "lucide-react";

export default function Dashboard() {
  return (
    <div className="min-h-screen bg-background">
      {/* Header */}
      <header className="border-b bg-card">
        <div className="container mx-auto px-4 py-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-3">
              <div className="rounded-lg bg-primary p-2">
                <Bot className="h-6 w-6 text-primary-foreground" />
              </div>
              <div>
                <h1 className="text-xl font-bold">Workflow Dashboard</h1>
                <p className="text-sm text-muted-foreground">
                  Manage and monitor your workflows
                </p>
              </div>
            </div>
            <CreateWorkflowModal />
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="container mx-auto px-4 py-8">
        {/* Workflows Section */}
        <section>
          <h2 className="text-lg font-semibold mb-4">My Workflows</h2>
          <WorkflowList />
        </section>

        {/* Schedule Section (Future Feature Placeholder) */}
        <ScheduleSection />
      </main>
    </div>
  );
}
