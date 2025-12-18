package esxi

import (
	"fmt"
	"log"
	"math/rand"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceVIRTUALDISKCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Config)
	//esxiConnInfo := ConnectionStruct{c.esxiHostName, c.esxiHostPort, c.esxiUserName, c.esxiPassword}
	log.Println("[resourceVIRTUALDISKCreate]")

	virtual_disk_disk_store := d.Get("virtual_disk_disk_store").(string)
	virtual_disk_dir := d.Get("virtual_disk_dir").(string)
	virtual_disk_name := d.Get("virtual_disk_name").(string)
	virtual_disk_size := d.Get("virtual_disk_size").(int)
	virtual_disk_type := d.Get("virtual_disk_type").(string)
	virtual_disk_clone_disk_store := d.Get("virtual_disk_clone_disk_store").(string)
	virtual_disk_clone_dir := d.Get("virtual_disk_clone_dir").(string)
	virtual_disk_clone_src_name := d.Get("virtual_disk_clone_src_name").(string)

	if virtual_disk_name == "" {
		const digits = "0123456789ABCDEF"
		name := make([]byte, 10)
		for i := range name {
			name[i] = digits[rand.Intn(len(digits))]
		}

		virtual_disk_name = fmt.Sprintf("vdisk_%s.vmdk", name)
	}

	//
	//  Validate virtual_disk_name
	//

	// todo,  check invalid chars (quotes, slash, period, comma)

	if !strings.HasSuffix(virtual_disk_name, ".vmdk") {
		return fmt.Errorf("virtual_disk_name does not have '.vmdk' suffix")
	}

	if !strings.HasSuffix(virtual_disk_clone_src_name, ".vmdk") {
		return fmt.Errorf("virtual_disk_clone_src_path does not have '.vmdk' suffix")
	}

	// Clone virtual disk

	if virtual_disk_clone_disk_store != "" || virtual_disk_clone_dir != "" && virtual_disk_clone_src_name != "" {

		virtdisk_id, err := virtualDiskCLONE(c, virtual_disk_disk_store, virtual_disk_dir,
			virtual_disk_name, virtual_disk_size, virtual_disk_type,
			virtual_disk_clone_disk_store,
			virtual_disk_clone_dir,
			virtual_disk_clone_src_name,
		)
		if err == nil {
			d.SetId(virtdisk_id)
		} else {
			log.Println("[resourceVIRTUALDISKClone] Error: " + err.Error())
			d.SetId("")
			return fmt.Errorf("failed to create virtual Disk: %s, Error: %s", virtual_disk_name, err.Error())
		}

	}

	// Create virtual disk

	virtdisk_id, err := virtualDiskCREATE(c, virtual_disk_disk_store, virtual_disk_dir,
		virtual_disk_name, virtual_disk_size, virtual_disk_type)
	if err == nil {
		d.SetId(virtdisk_id)
	} else {
		log.Println("[resourceVIRTUALDISKCreate] Error: " + err.Error())
		d.SetId("")
		return fmt.Errorf("failed to create virtual Disk: %s, Error: %s", virtual_disk_name, err.Error())
	}

	return nil
}
