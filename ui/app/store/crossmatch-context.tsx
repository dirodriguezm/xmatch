"use client";

import {
  createContext,
  Dispatch,
  ReactNode,
  useContext,
  useReducer,
} from "react";

// Types
export interface ResolverState {
  targetName: string;
  service: "SIMBAD" | "NED" | "VizieR";
  isResolving: boolean;
}

export type ResultsState = "empty" | "loading" | "success" | "error";

export interface CrossmatchState {
  resolver: ResolverState;
  resultsState: ResultsState;
}

export type CrossmatchAction =
  | { type: "SET_RESOLVER"; payload: Partial<ResolverState> }
  | { type: "SET_RESULTS_STATE"; payload: ResultsState }
  | { type: "RESET" };

const initialState: CrossmatchState = {
  resolver: {
    targetName: "",
    service: "SIMBAD",
    isResolving: false,
  },
  resultsState: "empty",
};

function crossmatchReducer(
  state: CrossmatchState,
  action: CrossmatchAction
): CrossmatchState {
  switch (action.type) {
    case "SET_RESOLVER":
      return { ...state, resolver: { ...state.resolver, ...action.payload } };
    case "SET_RESULTS_STATE":
      return { ...state, resultsState: action.payload };
    case "RESET":
      return initialState;
    default:
      return state;
  }
}

// Context
const CrossmatchContext = createContext<{
  state: CrossmatchState;
  dispatch: Dispatch<CrossmatchAction>;
} | null>(null);

export function CrossmatchProvider({ children }: { children: ReactNode }) {
  const [state, dispatch] = useReducer(crossmatchReducer, initialState);
  return (
    <CrossmatchContext.Provider value={{ state, dispatch }}>
      {children}
    </CrossmatchContext.Provider>
  );
}

export function useCrossmatchState() {
  const context = useContext(CrossmatchContext);
  if (!context) {
    throw new Error(
      "useCrossmatchState must be used within CrossmatchProvider"
    );
  }
  return context;
}
