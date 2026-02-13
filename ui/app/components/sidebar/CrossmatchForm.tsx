"use client";

import { SearchOutlined } from "@ant-design/icons";
import { Button } from "antd";

import { useSearchParams } from "@/app/hooks/useSearchParamsSync";
import { useCrossmatchState } from "@/app/store/crossmatch-context";

import { SearchParameters } from "./SearchParameters";
import { TargetCoordinates } from "./TargetCoordinates";

interface CrossmatchFormProps {
  onSubmit?: () => void;
}

export function CrossmatchForm({ onSubmit }: CrossmatchFormProps) {
  const { state, dispatch } = useCrossmatchState();
  const { isValid } = useSearchParams();

  const handleSubmit = () => {
    dispatch({ type: "SET_RESULTS_STATE", payload: "loading" });
    if (onSubmit) {
      onSubmit();
    }
  };

  return (
    <div className="flex flex-col h-full">
      <div className="flex-1 overflow-auto">
        <TargetCoordinates />
        <SearchParameters />
      </div>

      {/* Run Crossmatch Button - Fixed at bottom */}
      <div className="sidebar-section border-t border-border">
        <Button
          type="primary"
          size="large"
          block
          icon={<SearchOutlined />}
          onClick={handleSubmit}
          disabled={!isValid}
          loading={state.resultsState === "loading"}
        >
          Run Crossmatch
        </Button>
      </div>
    </div>
  );
}
