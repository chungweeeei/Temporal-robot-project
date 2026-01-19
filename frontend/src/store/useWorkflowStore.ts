import { create } from 'zustand';

interface WorkflowStore {
    activeWorkflowId: string | null;
    setActiveWorkflowId: (id: string | null) => void;
}

export const useWorkflowStore = create<WorkflowStore>((set) => ({
    activeWorkflowId: null,
    setActiveWorkflowId: (id: string | null) => set({ activeWorkflowId: id }),
}));