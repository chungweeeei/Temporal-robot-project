import { create } from 'zustand';

interface WorkflowState {
    activeWorkflowId: string | null;
    setActiveWorkflowId: (id: string | null) => void;
}

export const useWorkflowStore = create<WorkflowState>((set) => ({
    activeWorkflowId: null,
    setActiveWorkflowId: (id: string | null) => set({ activeWorkflowId: id }),
}));