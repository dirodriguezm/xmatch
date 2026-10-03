"use client";

import {
  ClearOutlined,
  ExperimentOutlined,
  UploadOutlined,
} from "@ant-design/icons";
import { Alert, App, Button, Input, Typography, Upload } from "antd";

import {
  BULK_EXAMPLE,
  type BulkParseResult,
  MAX_BULK_ROWS,
} from "@/app/lib/utils/bulkInput";

const { Text } = Typography;

/** Refuse files large enough to freeze the textarea; 1000 rows fit in ~100 kB. */
const MAX_FILE_BYTES = 2 * 1024 * 1024;
const MAX_ERRORS_SHOWN = 20;

interface BulkInputPanelProps {
  value: string;
  onChange: (value: string) => void;
  parsed: BulkParseResult;
}

export function BulkInputPanel({
  value,
  onChange,
  parsed,
}: BulkInputPanelProps) {
  const { message } = App.useApp();
  const { targets, errors, header } = parsed;
  const tooMany = targets.length > MAX_BULK_ROWS;

  const handleFile = (file: File) => {
    if (file.size > MAX_FILE_BYTES) {
      message.error(
        `${file.name} is ${(file.size / 1024 / 1024).toFixed(1)} MB; the limit is 2 MB (${MAX_BULK_ROWS} rows).`
      );
      return false;
    }
    file
      .text()
      .then((text) => {
        onChange(text);
        message.success(`Loaded ${file.name}`);
      })
      .catch(() => message.error(`Could not read ${file.name}`));
    // Keep antd from uploading anywhere; we only need the text.
    return false;
  };

  return (
    <div className="flex flex-col gap-3">
      <div className="flex flex-wrap items-center justify-between gap-2">
        <Text strong className="text-foreground">
          Positions
        </Text>
        <div className="flex flex-wrap gap-2">
          <Upload
            accept=".csv,.tsv,.txt,.dat,text/csv,text/plain,text/tab-separated-values"
            beforeUpload={handleFile}
            showUploadList={false}
            maxCount={1}
          >
            <Button size="small" icon={<UploadOutlined />}>
              Upload CSV/TSV
            </Button>
          </Upload>
          <Button
            size="small"
            icon={<ExperimentOutlined />}
            onClick={() => onChange(BULK_EXAMPLE)}
          >
            Load example
          </Button>
          <Button
            size="small"
            icon={<ClearOutlined />}
            disabled={!value}
            onClick={() => onChange("")}
          >
            Clear
          </Button>
        </div>
      </div>

      <Input.TextArea
        value={value}
        onChange={(e) => onChange(e.target.value)}
        placeholder={
          "name,ra,dec\nM31,10.684708,41.268750\nSirius,06:45:08.917,-16:42:58.02"
        }
        autoSize={{ minRows: 8, maxRows: 18 }}
        spellCheck={false}
        className="!font-mono !text-xs"
        aria-label="Positions to cross-match"
      />

      <Text className="text-neutral-400 text-xs">
        One position per line. Comma, tab, semicolon or space separated; RA/Dec
        in decimal degrees or sexagesimal (hh:mm:ss ±dd:mm:ss). An optional
        header with <code>ra</code>, <code>dec</code> and <code>name</code>{" "}
        columns is detected automatically. Lines starting with # are ignored. Up
        to {MAX_BULK_ROWS} rows.
      </Text>

      {value.trim() && (
        <div className="flex flex-wrap gap-x-4 gap-y-1 text-xs">
          <span className="text-foreground">
            {targets.length} position{targets.length === 1 ? "" : "s"} parsed
          </span>
          {header && (
            <span className="text-neutral-400">
              header: {header.join(", ")}
            </span>
          )}
          {errors.length > 0 && (
            <span className="text-red-400">
              {errors.length} line{errors.length === 1 ? "" : "s"} skipped
            </span>
          )}
        </div>
      )}

      {tooMany && (
        <Alert
          type="error"
          showIcon
          title={`Too many positions: ${targets.length}. The limit is ${MAX_BULK_ROWS} per run — split the list, or use the API directly for larger jobs.`}
        />
      )}

      {errors.length > 0 && (
        <Alert
          type="warning"
          showIcon
          title="Some lines could not be parsed and will be skipped"
          description={
            <ul className="m-0 pl-4 text-xs space-y-0.5 max-h-48 overflow-auto">
              {errors.slice(0, MAX_ERRORS_SHOWN).map((e) => (
                <li key={e.line}>
                  <span className="font-mono text-neutral-400">
                    line {e.line}:
                  </span>{" "}
                  <span className="font-mono">
                    {e.text.length > 60 ? `${e.text.slice(0, 60)}…` : e.text}
                  </span>{" "}
                  — {e.message}
                </li>
              ))}
              {errors.length > MAX_ERRORS_SHOWN && (
                <li className="text-neutral-400">
                  …and {errors.length - MAX_ERRORS_SHOWN} more
                </li>
              )}
            </ul>
          }
        />
      )}
    </div>
  );
}
