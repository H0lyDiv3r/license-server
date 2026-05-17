import { Button } from "@/components/ui/button";
import { WriteJournalEntry } from "../../../wailsjs/go/license/License";

export const MainPage = () => {
  const handleComputeHmac = () => {
    // WriteJournalEntry();
  };

  return (
    <>
      main app
      <Button
        onClick={() => {
          handleComputeHmac();
        }}
      >
        compute hmac
      </Button>
    </>
  );
};
