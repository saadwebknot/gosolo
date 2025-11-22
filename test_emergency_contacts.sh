#!/bin/bash
# Test script for updating emergency contacts

echo "Testing Update Emergency Contacts API..."

curl -X PUT http://localhost:8005/gosolo/travellers/1/emergency-contacts \
  -H "Content-Type: application/json" \
  -d '{"emergency_contact_phone":"+1-202-555-0123"}'

echo ""
echo ""
echo "Done!"

