import React from "react";
import { useNavigate } from "react-router-dom";
import { Button, InlineNotification } from '@carbon/react';

function NotFoundPage() {
  const navigate = useNavigate();
  return (
    <div style={{ display: "flex", flexDirection: "column", alignItems: "center", marginTop: "6rem" }}>
      <InlineNotification
        kind="error"
        title="404 - Page Not Found"
        subtitle="The page you are looking for does not exist."
        hideCloseButton
        style={{ maxWidth: "400px", marginBottom: "2rem" }}
      />
      <Button kind="primary" size="lg" onClick={() => navigate("/")}>
        Go to Home
      </Button>
    </div>
  );
}

export default NotFoundPage;
